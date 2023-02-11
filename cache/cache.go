package cache

import (
	"encoding/json"
	"fmt"
)

var instance Cache

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	close()
}

func Setup() {
	instance = newInMemoryCache(newMemoryParams{
		persist:  true,
		fileName: ".cache.json",
	})
}

func Close() {
	instance.close()
}

func Get[T any](key string, res *T) error {
	fmt.Println("asking for", key)
	data, err := instance.Get(key)
	if err != nil {
		return err
	}

	bytes, ok := data.([]byte)
	if !ok {
		str := data.(string)
		bytes = []byte(str)
	}

	err = json.Unmarshal(bytes, res)
	if err != nil {
		return err
	}

	return nil
}

func Set(key string, value interface{}) error {
	return instance.Set(key, value)
}
