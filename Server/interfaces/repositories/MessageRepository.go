package repositories

import (
	"context"
	"redisChat/Server/repositories/entities"
)

type MessageRepository interface {

	// Flush is responsible for flushing data from the repo
	Flush(ctx context.Context)

	// Save is reponsible for saving a new entry in the repo
	Save(ctx context.Context, message *entities.Message) (err error)

	// RetrieveByMessageID is reponsible for getting an existing entry from the repo
	RetrieveByMessageID(ctx context.Context, messageID int) (resultMessage *entities.Message, err error)

	// Retrieve is reponsible for getting an existing entry from the repo
	RetrieveAll(ctx context.Context) (resultSlice *entities.Messages, err error)

	// Delete is reponsible for deleting an entry from the repo
	Delete(ctx context.Context, messageID int) (err error)
}
