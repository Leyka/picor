package cache

import "encoding/json"

var instance Cache

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	close()
}

func Setup() {
	instance = newInMemoryCache(newMemoryParams{
		persist:  true,
		fileName: ".cache",
	})

	instance.Set("test", "test")
}

func Close() {
	instance.close()
}

func Get[T any](key string, res T) error {
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
