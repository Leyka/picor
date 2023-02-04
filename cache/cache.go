package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Instance Cache

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

func SetupCache() *Cache {
	// Check if redis is running otherwise use in memory cache
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		Instance = newInMemoryCache()
	} else {
		Instance = newRedisCache(redisClient, context.Background())
	}

	return &Instance
}
