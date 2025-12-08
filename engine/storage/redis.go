package storage

import (
	"context"

	"github.com/acer-red/home/engine/sys"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// InitRedis initializes the shared Redis client.
func InitRedis(cfg sys.Redis) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb.Ping(context.Background()).Err()
}

// CloseRedis closes the Redis client if present.
func CloseRedis() error {
	if rdb == nil {
		return nil
	}
	return rdb.Close()
}

// redisClient returns the initialized client or nil.
func redisClient() *redis.Client {
	return rdb
}
