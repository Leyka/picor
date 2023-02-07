package cache

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
