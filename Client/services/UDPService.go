package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"redisChat/Client/config"
	"redisChat/Client/interfaces/repositories"
	"redisChat/Client/interfaces/services"
	logger "redisChat/Client/pkg/log"
	"redisChat/Client/repositories/entities"
	"redisChat/Client/services/viewmodels"
	"strings"
)

type udpService struct {
	messageRepository repositories.MessageRepository
}

func NewUDPService(messageRepo repositories.MessageRepository) services.UDPService {
	return &udpService{
		messageRepository: messageRepo,
	}
}

// ReceiveMessage unmarshals a message received via udp and calls the handleMessage function
func (s *udpService) ReceiveMessage(buf []byte, n int) {

	// Check if message is an error
	if strings.Contains(string(buf[:n]), "err") {
		logger.Logger.Info("received an error message", logger.Information{
			"error": string(buf[:n]),
		})
		return
	}

	// Check if server is closing
	if strings.Contains(string(buf[:n]), "closing connection") {
		logger.Logger.Info("server has closed", logger.Information{})
		os.Exit(1)
	}

	// Unmarshal the byte slice
	udpMessages := viewmodels.UDPMessages{}
	var err error
	if err = json.Unmarshal(buf[:n], &udpMessages); err != nil {
		logger.Logger.Error("unmarshaling message", err, logger.Information{
			"buf bytes": buf,
		})
		return
	}

	// Handle the message
	s.handleMessage(&udpMessages)
}

// SendMessages marshals the provided message and sends it via udp
func (s *udpService) SendMessage(message *viewmodels.UDPMessage) (err error) {
	// Get server address and port from config
	address := config.ClientConfig.ServerAddressReceiver
	myPort := config.ClientPort

	// Initiate UDP to send messages
	sender, err := net.ListenPacket("udp", fmt.Sprintf(":%s", myPort))
	if err != nil {
		logger.Logger.Error("starting up udp sender", err, logger.Information{
			"server port": myPort,
		})
		return
	}
	defer sender.Close()

	// Resolve send address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		logger.Logger.Error("resolving address for udp broadcast", err, logger.Information{
			"address": address,
		})
	}

	// Marshal message
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Error("marshaling messages", err, logger.Information{
			"message": message,
		})
		return
	}

	// Send message
	_, err = sender.WriteTo(msg, addr)
	if err != nil {
		logger.Logger.Error("sending udp mesage", err, logger.Information{
			"message": message,
			"address": address,
		})
	}

	// Get response
	buf := make([]byte, 1024)
	n, _, err := sender.ReadFrom(buf)
	if err != nil {
		logger.Logger.Error("reading udp response message", err, logger.Information{})
	}

	if strings.Contains(string(buf[:n]), "err") {
		logger.Logger.Info("received an error message", logger.Information{})
		err = errors.New(string(buf[:n]))
	}

	return
}

// DeleteMessage marshals the delete message and sends it via udp
func (s *udpService) DeleteMessage(messageID int, message *viewmodels.UDPMessage) (err error) {

	// Check that user is owner of message
	messageUser := s.messageRepository.GetMessageUserByID(messageID)
	if message.User != messageUser {
		err = services.ErrPermissionDenied
		return
	}

	// Send message to delete message
	if err = s.SendMessage(&viewmodels.UDPMessage{
		Message: fmt.Sprintf("delete %d", messageID),
		User:    config.User,
	}); err != nil {
		logger.Logger.Error("sending udp delete mesage", err, logger.Information{})
		return
	}

	return
}

// GetAllMessages retrieves all messages that are stored in memory
// Set on connect and updated with every broadcast
func (s *udpService) GetAllMessages() (messages *viewmodels.UDPMessages) {
	entityMessages := s.messageRepository.GetMessages()
	messages = &viewmodels.UDPMessages{}
	entityMessages.ToViewmodel(messages)
	return
}

// handleMessage converts to the correct type and calls required repo function
func (s *udpService) handleMessage(messages *viewmodels.UDPMessages) {

	// Convert from viewmodel to entity
	entityMessages := &entities.Messages{}
	entityMessages.ToEntity(messages)

	// Call repo function that will update map
	s.messageRepository.UpdateMessages(entityMessages)
}
