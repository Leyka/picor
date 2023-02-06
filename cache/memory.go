package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/Leyka/picor/file"
)

// Used to dump the cache to a file when the program exits
const DATA_FILE = ".cache.json"

type inMemoryCache struct {
	data    *map[string]interface{}
	persist bool
	mu      sync.Mutex
}

func newInMemoryCache(persist bool) *inMemoryCache {
	var m map[string]interface{}

	if file.Exists(DATA_FILE) {
		m = *load()
	} else {
		m = make(map[string]interface{})
	}

	return &inMemoryCache{
		data:    &m,
		persist: persist,
	}
}

func (c *inMemoryCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value, ok := (*c.data)[key]; ok {
		return value, nil
	}

	err := fmt.Errorf("key %s not found", key)
	return nil, err
}

func (c *inMemoryCache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	(*c.data)[key] = value
	return nil
}

func (c *inMemoryCache) close() {
	if len(*c.data) == 0 {
		return
	}

	if c.persist {
		c.dump()
	}
}

func (c *inMemoryCache) dump() {
	// Deserialize the cache (map) to a json file
	b, err := json.Marshal(*c.data)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = ioutil.WriteFile(DATA_FILE, b, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func load() *map[string]interface{} {
	// Load the cache from a file
	b, err := ioutil.ReadFile(DATA_FILE)
	if err != nil {
		fmt.Println("error:", err)
		return &map[string]interface{}{}
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("error:", err)
		return &map[string]interface{}{}
	}

	return &m
}
