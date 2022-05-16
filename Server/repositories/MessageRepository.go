package repositories

import (
	"context"
	"redisChat/Server/interfaces/repositories"
	logger "redisChat/Server/pkg/log"
	"redisChat/Server/pkg/redisManager"
	"redisChat/Server/repositories/entities"
	"strconv"
)

type messageRepository struct {
}

func NewMessageRepository() repositories.MessageRepository {
	return &messageRepository{}
}

// Flush is responsible for flushing data from the repo
func (r *messageRepository) Flush(ctx context.Context) {
	redisManager.Client.Flush(ctx)
}

// Save is reponsible for saving a new entry in the repo
func (r *messageRepository) Save(ctx context.Context, message *entities.Message) (err error) {

	// Get last ID
	lastID, err := r.getLastID(ctx)
	if err != nil {
		return
	}

	// Add ID to message
	message.ID = lastID

	// Save to the db/cache
	err = redisManager.Client.Save(ctx, lastID, message)
	if err != nil {
		return
	}

	// Increment lastID
	err = r.incrementID(ctx, lastID)
	return
}

// RetrieveByKey is reponsible for getting an existing entry from the redis db/cache by key
func (r *messageRepository) RetrieveByMessageID(ctx context.Context, messageID int) (resultMessage *entities.Message, err error) {
	result := entities.Message{}
	if err = redisManager.Client.RetrieveByKey(ctx, messageID, &result); err != nil {
		logger.Logger.Error("retrieving message by messageID", err, logger.Information{})
	}

	resultMessage = &result
	return
}

// Retrieve is reponsible for getting an existing entry from the redis db/cache
func (r *messageRepository) RetrieveAll(ctx context.Context) (resultMessages *entities.Messages, err error) {
	// Initiate variables
	resultMessages = &entities.Messages{}
	messages := []entities.Message{}
	keys, err := redisManager.Client.RetrieveAll(ctx)
	if err != nil {
		logger.Logger.Error("getting all message keys", err, logger.Information{})
	}

	// Loop through all keys to get value
	for _, key := range keys {
		// Convert key to int
		var k int
		var message *entities.Message
		k, err = strconv.Atoi(key)
		if err != nil {
			return
		}

		// Get message
		message, err = r.RetrieveByMessageID(ctx, k)
		if err != nil {
			return
		}

		// Append to slice
		messages = append(messages, *message)
	}

	resultMessages.Messages = messages
	return
}

// Delete is reponsible for deleting an entry from the repo
func (r *messageRepository) Delete(ctx context.Context, messageID int) (err error) {
	if err = redisManager.Client.Delete(ctx, messageID); err != nil {
		logger.Logger.Error("retrieving message by messageID", err, logger.Information{
			"message ID": messageID,
		})
	}
	return
}

// incrementID increments the value of ID and calls the db function
func (r *messageRepository) incrementID(ctx context.Context, ID int) (err error) {
	err = redisManager.Client.SetLastID(ctx, "ID", ID+1)
	if err != nil {
		logger.Logger.Error("incrementing last id", err, logger.Information{})
	}
	return
}

// getLastID returns the last used ID
func (r *messageRepository) getLastID(ctx context.Context) (ID int, err error) {
	// Get ID
	ID, err = redisManager.Client.GetLastID(ctx, "ID")
	if err != nil {
		logger.Logger.Error("getting last id", err, logger.Information{})
	}
	return
}
