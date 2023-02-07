package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/Leyka/picor/file"
)

type inMemoryCache struct {
	data     *map[string]interface{}
	mu       sync.Mutex
	persist  bool
	fileName string
}

type newMemoryParams struct {
	persist  bool
	fileName string
}

func newInMemoryCache(params newMemoryParams) *inMemoryCache {
	var data map[string]interface{}

	if file.Exists(params.fileName) {
		data, err := deserialize(params.fileName)
		if err == nil {
			return &inMemoryCache{
				data:     data,
				persist:  params.persist,
				fileName: params.fileName,
			}
		}
	}

	data = make(map[string]interface{})
	return &inMemoryCache{
		data:     &data,
		persist:  params.persist,
		fileName: params.fileName,
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
	if len(*c.data) == 0 || !c.persist {
		return
	}

	c.serialize()
}

func (c *inMemoryCache) serialize() error {
	b, err := json.Marshal(*c.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.fileName, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func deserialize(fileName string) (*map[string]interface{}, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
