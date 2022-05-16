package redisManager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"redisChat/Server/config"
	"strconv"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

var ErrNotFound = errors.New("entry could not be found")

type RedisClient struct {
	client *redis.Client
	cache  *cache.Cache
}

// NewRedisClient creates a new client with the required configurations
func NewRedisClient(address, password string, db int) Redis {
	// Define redis db
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	// Define cache
	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(20, time.Hour*24),
	})
	return &RedisClient{
		client: client,
		cache:  cache,
	}
}

// SendPingMessage sends a ping message and returns a pong message
func (c *RedisClient) SendPingMessage(ctx context.Context) (string, error) {
	return c.client.Ping(ctx).Result()
}

// Close closes the connection to the redis db
func (c *RedisClient) Close() error {
	return c.client.Close()
}

// Flush is responsible for flushing the redis db
func (c *RedisClient) Flush(ctx context.Context) {
	c.client.FlushDB(ctx)
}

// SetLastID is reponsible for setting the last ID in redis db/cache
func (c *RedisClient) SetLastID(ctx context.Context, key string, value int) (err error) {

	// Set value in redis
	err = c.client.Set(ctx, key, fmt.Sprintf("%d", value), 0).Err()
	return
}

// GetLastID is reponsible for getting the last ID in redis db/cache
func (c *RedisClient) GetLastID(ctx context.Context, key string) (value int, err error) {

	// Get last ID
	var redisVal string
	redisVal, err = c.client.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return
		}
		err = nil
	}

	// If ID has not been set before, set it to 1 now.
	if redisVal == "" {
		err = c.SetLastID(ctx, "ID", 1)
		if err != nil {
			return
		}

		value = 1
		return
	}

	// Convert ID to int
	value, err = strconv.Atoi(redisVal)
	return
}

// Save is reponsible for saving a new entry in the redis cache
func (c *RedisClient) Save(ctx context.Context, key int, value interface{}) (err error) {

	// Marshal value
	val, err := json.Marshal(value)
	if err != nil {
		return
	}

	// Set value in redis
	err = c.client.Set(ctx, fmt.Sprintf("%d", key), val, 0).Err()
	if err != nil {
		return
	}

	// Check that redis limit has not been reached
	err = c.limit(ctx)
	return
}

// RetrieveByKey is reponsible for getting an existing entry from the redis cache by key
func (c *RedisClient) RetrieveByKey(ctx context.Context, key int, value interface{}) (err error) {
	val, err := c.client.Get(ctx, fmt.Sprintf("%d", key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = ErrNotFound
		}
	}

	json.Unmarshal([]byte(val), value)
	return
}

// Retrieve is reponsible for getting an existing entry from the redis db
func (c *RedisClient) RetrieveAll(ctx context.Context) (keys []string, err error) {
	iter := c.client.Scan(ctx, 0, "*", 0).Iterator()
	if err = iter.Err(); err != nil {
		return
	}

	for iter.Next(ctx) {
		key := iter.Val()

		if key != "ID" {
			keys = append(keys, key)
		}
	}

	return
}

// Retrieve is reponsible for getting an existing entry from the redis db
func (c *RedisClient) limit(ctx context.Context) (err error) {

	// Scan for all keys
	count := 0
	IDs := []int{}
	iter := c.client.Scan(ctx, 0, "*", 0).Iterator()
	if err = iter.Err(); err != nil {
		return
	}

	// Count keys and save to slice
	for iter.Next(ctx) {
		key := iter.Val()

		if key != "ID" {
			var ID int
			ID, err = strconv.Atoi(key)
			if err != nil {
				return
			}
			IDs = append(IDs, ID)
		}
		count++
	}

	// If there are more than max entries, delete oldest
	if count > config.ServerConfig.MaxEntries {

		// Get oldest ID
		oldestID := min(IDs)

		// Delete oldest ID
		err = c.Delete(ctx, oldestID)
	}

	return
}

// Delete is reponsible for deleting an entry from the redis db/cache
func (c *RedisClient) Delete(ctx context.Context, key int) (err error) {
	err = c.client.Del(ctx, fmt.Sprintf("%d", key)).Err()
	return
}

func min(slice []int) (min int) {
	min = math.MaxInt64
	for _, s := range slice {
		if s < min {
			min = s
		}
	}
	return
}
