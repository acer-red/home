package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// initializes the shared Redis client.
func Init(opt redis.Options) error {
	rdb = redis.NewClient(&opt)
	return rdb.Ping(context.Background()).Err()
}

// closes the Redis client if present.
func Close() error {
	if rdb == nil {
		return nil
	}
	return rdb.Close()
}

// redisClient returns the initialized client or nil.
func redisClient() *redis.Client {
	return rdb
}
