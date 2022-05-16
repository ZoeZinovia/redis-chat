package repositories

import (
	"redisChat/Client/interfaces/repositories"
	"redisChat/Client/repositories/entities"
	"sort"
	"sync"
)

var mapMutex sync.Mutex

type messageRepository struct {
	messageMap map[int]entities.Message
}

func NewMessageRepository() repositories.MessageRepository {
	mMap := make(map[int]entities.Message)
	return &messageRepository{
		messageMap: mMap,
	}
}

// UpdateMessages overwrites existing map with new map
func (r *messageRepository) UpdateMessages(messages *entities.Messages) {
	mapMutex.Lock()
	newMap := make(map[int]entities.Message)
	for _, msg := range *messages.Messages {
		newMap[msg.ID] = msg
	}

	// Overwrite old map
	r.messageMap = newMap
	mapMutex.Unlock()
}

// GetMessages gets all messages from the map
func (r *messageRepository) GetMessages() (messages *entities.Messages) {
	mapMutex.Lock()
	msgs := []entities.Message{}
	for _, msg := range r.messageMap {
		msgs = append(msgs, msg)
	}
	mapMutex.Unlock()

	// Sort messages by id
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID > msgs[j].ID
	})

	// Return messages
	messages = &entities.Messages{
		Messages: &msgs,
	}

	return
}
