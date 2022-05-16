package redisManager

import "context"

type Redis interface {
	// SendPingMessage sends a ping message and returns a pong message
	SendPingMessage(context.Context) (string, error)

	// Close closes the connection to the redis db
	Close() error

	// Flush is responsible for flushing the redis db
	Flush(ctx context.Context)

	// SetLastID is reponsible for setting the last ID in redis db/cache
	SetLastID(ctx context.Context, key string, value int) (err error)

	// SetLastID is reponsible for setting the last ID in redis db/cache
	GetLastID(ctx context.Context, key string) (value int, err error)

	// Save is reponsible for saving a new entry in the redis db/cache
	Save(ctx context.Context, key int, value interface{}) (err error)

	// RetrieveByKey is reponsible for getting an existing entry from the redis db/cache by key
	RetrieveByKey(ctx context.Context, key int, value interface{}) (err error)

	// Retrieve is reponsible for getting an existing entry from the redis db/cache
	RetrieveAll(ctx context.Context) (keys []string, err error)

	// Delete is reponsible for deleting an entry from the redis db/cache
	Delete(ctx context.Context, key int) (err error)
}

var Client Redis = NewRedisClient("localhost:6379", "", 0)
