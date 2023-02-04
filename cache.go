package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

func InitCache() Cache {
	// Check if redis is running otherwise use in memory cache
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return newInMemoryCache()
	}

	return newRedisCache(redisClient, context.Background())
}

// ~ In memory caching ~
type inMemoryCache struct {
	data map[string]interface{}
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{data: make(map[string]interface{})}
}

func (c *inMemoryCache) Get(key string) (interface{}, error) {
	return c.data[key], nil
}

func (c *inMemoryCache) Set(key string, value interface{}) error {
	c.data[key] = value
	return nil
}

// ~ Redis caching ~
type redisCache struct {
	client  *redis.Client
	context context.Context
}

func newRedisCache(client *redis.Client, context context.Context) *redisCache {
	return &redisCache{
		client:  client,
		context: context,
	}
}
func (c *redisCache) Get(key string) (interface{}, error) {
	return c.client.Get(c.context, key).Result()
}
func (c *redisCache) Set(key string, value interface{}) error {
	return c.client.Set(c.context, key, value, 0).Err()
}
