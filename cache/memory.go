package cache

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
