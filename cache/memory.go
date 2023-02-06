package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"

	"github.com/Leyka/picor/file"
)

// Used to dump the cache to a file when the program exits
const DATA_FILE = ".data.json"

type inMemoryCache struct {
	data    sync.Map
	count   *atomic.Uint64
	persist bool
}

func newInMemoryCache(persist bool) *inMemoryCache {
	var m sync.Map

	if file.Exists(DATA_FILE) {
		m = *load()
	}

	return &inMemoryCache{
		data:    m,
		count:   new(atomic.Uint64),
		persist: persist,
	}
}

func (c *inMemoryCache) Get(key string) (interface{}, error) {
	if value, ok := c.data.Load(key); ok {
		return value, nil
	}

	err := fmt.Errorf("key %s not found", key)
	return nil, err
}

func (c *inMemoryCache) Set(key string, value interface{}) error {
	if _, ok := c.data.Load(key); !ok {
		// Increment the counter if new key is added
		c.count.Add(1)
	}

	c.data.Store(key, value)
	return nil
}

func (c *inMemoryCache) close() {
	if c.count == nil || c.count.Load() == 0 {
		return
	}

	if c.persist {
		c.dump()
	}
}

func (c *inMemoryCache) dump() {
	// Deserialize the cache (sync.Map) to a json file
	data := make(map[string]interface{})
	c.data.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = ioutil.WriteFile(DATA_FILE, b, 0644)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

func load() *sync.Map {
	// Load the cache from a file
	b, err := ioutil.ReadFile(DATA_FILE)
	if err != nil {
		fmt.Println("error:", err)
		return &sync.Map{}
	}

	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Println("error:", err)
		return &sync.Map{}
	}

	m := &sync.Map{}
	for key, value := range data {
		m.Store(key, value)
	}

	return m
}
