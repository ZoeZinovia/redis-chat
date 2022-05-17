package repositories

import "redisChat/Client/repositories/entities"

type MessageRepository interface {
	// UpdateMessages overwrites existing map with new map
	UpdateMessages(messages *entities.Messages)

	// GetMessages gets all messages from the map
	GetMessages() (messages *entities.Messages)

	// GetMessageUserByID gets the user associated with a message
	GetMessageUserByID(ID int) (messageUser string)
}
