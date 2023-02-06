package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var instance Cache

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	close()
}

func Setup() {
	// Redis running? If not then use in-memory cache
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		persist := true
		instance = newInMemoryCache(persist)
		return
	}

	instance = newRedisCache(redisClient, context.Background())
}

func Close() {
	instance.close()
}
