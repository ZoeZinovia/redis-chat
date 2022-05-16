package services

import (
	"context"
	"encoding/json"
	"fmt"
	"redisChat/Server/interfaces/repositories"
	"redisChat/Server/interfaces/services"
	logger "redisChat/Server/pkg/log"
	"redisChat/Server/repositories/entities"
	"redisChat/Server/services/viewmodels"
	"redisChat/Server/worker"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var udpMutex sync.Mutex

type udpService struct {
	messageRepository repositories.MessageRepository
	udpClientMap      map[string]string // map containing user (key) and ip address (value)
}

func NewUDPService(messageRepo repositories.MessageRepository) services.UDPService {
	clientMap := make(map[string]string)
	return &udpService{
		udpClientMap:      clientMap,
		messageRepository: messageRepo,
	}
}

// ReceiveMessage unmarshals a message received via udp and calls the handleMessage function
func (s *udpService) ReceiveMessage(buf []byte, address string, n int) (response string, err error) {

	// Unmarshal the byte slice
	udpMessage := viewmodels.UDPMessage{}
	if err = json.Unmarshal(buf[:n], &udpMessage); err != nil {
		logger.Logger.Error("unmarshaling message", err, logger.Information{
			"buf bytes": buf,
			"address":   address,
		})
		return
	}

	// Handle the message
	response, err = s.handleMessage(&udpMessage, address)
	if err != nil {
		logger.Logger.Error("handling message", err, logger.Information{
			"udp message": udpMessage,
			"address":     address,
		})
	}

	return
}

// handleMessage checks what type of message was sent an handles it appropriately
// A connect message adds the client to the client map and retrieves the last 20 messages
// A diconnect message simply removes the client from the client map
// A delete message calls the repo function that deletes a message and broadcasts the information
// Any other message is considered a new message and is added to the db/cache and is broadcasted
// to all clients
func (s *udpService) handleMessage(udpMessage *viewmodels.UDPMessage, address string) (response string, err error) {

	// Initiate context
	ctx := context.Background()

	// Check if client wants to connect, disconnect or delete a message
	if udpMessage.Message == "connect" {

		// Connect the client
		s.connect(udpMessage.User, address)
		fmt.Printf("=== Connected %s %s! ===\n", udpMessage.User, address)

		// Get last messages
		response, err = s.GetLastMessagesForUDP(ctx)
		return

	} else if udpMessage.Message == "disconnect" {

		// Disconnect the client
		s.disconnect(ctx, address)
		fmt.Printf("=== Disconnected %s %s! ===\n", udpMessage.User, address)

		return
	}

	// Client has sent a message. Check if the client is already connected
	if !s.isConnected(udpMessage.User) {
		err = services.ErrNotConnected
		return
	}

	// Update IP address if it's different than current address
	if val := s.udpClientMap[udpMessage.User]; val != address {
		s.connect(udpMessage.User, address)
	}

	// Client is connected. Check if it is a new message or a request for deletion
	if strings.Contains(udpMessage.Message, "delete") {

		// Split message
		deleteSlice := strings.Split(udpMessage.Message, " ")
		if len(deleteSlice) != 2 || deleteSlice[0] != "delete" {
			err = services.ErrBadRequest
			return
		}

		// Get ID of message to be deleted
		var messageID int
		messageID, err = strconv.Atoi(deleteSlice[1])
		if err != nil {
			logger.Logger.Error("getting delete message id", err, logger.Information{
				"delete message": udpMessage.Message,
			})
			err = services.ErrBadRequest
			return
		}

		// Get message from the repo to check if it is the same user
		var messageToBeDeleted *entities.Message
		messageToBeDeleted, err = s.messageRepository.RetrieveByMessageID(ctx, messageID)
		if err != nil {
			return
		}
		if udpMessage.User != messageToBeDeleted.User {
			err = services.ErrNoPermission
			logger.Logger.Error("deleting message", err, logger.Information{
				"message id": messageID,
			})
			return
		}

		// Call repo function that will delete the message
		err = s.messageRepository.Delete(ctx, messageID)
		if err != nil {
			logger.Logger.Error("deleting message", err, logger.Information{
				"message id": messageID,
			})
			return
		}

		// Get last messages after delete
		response, err = s.GetLastMessagesForUDP(ctx)
		if err != nil {
			logger.Logger.Error("getting last messages after deleting message", err, logger.Information{
				"message id": messageID,
			})
		}

		for _, clientAddr := range s.udpClientMap {
			worker.BroadcastChannel <- viewmodels.UDPBroadcastMessage{
				DestinationAddress: clientAddr,
				Message:            response,
			}
		}

		fmt.Println("=== Message deleted and broadcasted ===")

		return
	}

	// Call repo function that will add message to the db
	messageEntity := entities.ToEntity(udpMessage)

	if err = s.messageRepository.Save(ctx, messageEntity); err != nil {
		logger.Logger.Error("saving message", err, logger.Information{
			"message": messageEntity,
		})
		return
	}

	// Get last messages after add
	response, err = s.GetLastMessagesForUDP(ctx)
	if err != nil {
		logger.Logger.Error("getting last messages after addind message", err, logger.Information{
			"message": udpMessage.Message,
		})
	}

	// Broadcast message to all other clients
	for _, clientAddr := range s.udpClientMap {
		worker.BroadcastChannel <- viewmodels.UDPBroadcastMessage{
			DestinationAddress: clientAddr,
			Message:            response,
		}
	}

	fmt.Println("=== Done broadcasting message ===")
	return
}

// connect is responsible for adding an ip address to the UDP client map.
// Mutex is used to ensure that we don't write to the map simultaneously
// with different routines
func (s *udpService) connect(user string, address string) {
	udpMutex.Lock()
	s.udpClientMap[user] = address
	udpMutex.Unlock()
}

// disconnect is responsible for removing an ip address from the UDP client map.
// Mutex is used to ensure that we don't write to the map simultaneously
// with different routines
func (s *udpService) disconnect(ctx context.Context, user string) {
	udpMutex.Lock()
	delete(s.udpClientMap, user)
	udpMutex.Unlock()

	// Check if the last client just disconnected
	if len(s.udpClientMap) == 0 {

		// Flush redis db/cache
		s.messageRepository.Flush(ctx)
		fmt.Println("=== Flushed Redis ===")
	}
}

// isConnected is responsible for checking if a client already exists in the client map
func (s *udpService) isConnected(address string) (connected bool) {
	_, ok := s.udpClientMap[address]
	return ok
}

// getLastMessagesForUDP gets the messages stored in the db
func (s *udpService) GetLastMessagesForUDP(ctx context.Context) (response string, err error) {
	// Get last messages from the repo layer
	allMessages := &entities.Messages{}
	allMessages, err = s.messageRepository.RetrieveAll(ctx)
	if err != nil {
		logger.Logger.Error("getting all messages", err, logger.Information{})
		return
	}

	// Sort messages by id
	sort.Slice(allMessages.Messages, func(i, j int) bool {
		return allMessages.Messages[i].ID > allMessages.Messages[j].ID
	})

	// Marshal to json and return
	var result []byte
	result, err = json.Marshal(allMessages)
	if err != nil {
		logger.Logger.Error("marshaling messages", err, logger.Information{
			"messages": allMessages,
		})
		return
	}

	response = string(result)
	return
}

// SendCloseMessages sends a close message to all clients that are still connected
func (s *udpService) SendCloseMessages() {
	for _, clientAddr := range s.udpClientMap {
		worker.BroadcastChannel <- viewmodels.UDPBroadcastMessage{
			DestinationAddress: clientAddr,
			Message:            "closing connection",
		}
	}
}
