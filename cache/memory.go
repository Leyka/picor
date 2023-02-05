package cache

import (
	"fmt"
	"sync"
)

type inMemoryCache struct {
	data sync.Map
}

func newInMemoryCache() *inMemoryCache {
	var m sync.Map
	return &inMemoryCache{data: m}
}

func (c *inMemoryCache) Get(key string) (interface{}, error) {
	if value, ok := c.data.Load(key); ok {
		return value, nil
	}

	err := fmt.Errorf("key %s not found", key)
	return nil, err
}

func (c *inMemoryCache) Set(key string, value interface{}) error {
	c.data.Store(key, value)
	return nil
}

func (c *inMemoryCache) close() {
	c.data = sync.Map{}
}
