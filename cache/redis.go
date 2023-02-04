package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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
