package cache

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

// Cache represents simple in-memory key/value store that can be persisted
type Cache struct {
	path    string
	data    map[string]string
	updated bool
	fake    bool
}

// New instantiates new Cache
func New(name string, notFake bool) (*Cache, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	home := user.HomeDir
	p := path.Join(home, ".cache", name+".json")

	cache := &Cache{
		path: p,
		data: make(map[string]string),
		fake: !notFake,
	}

	if notFake {
		f, err := os.ReadFile(cache.path)
		if err == nil {
			err = json.Unmarshal(f, &cache.data)
			if err != nil {
				return nil, err
			}
		}
	}

	return cache, nil
}

// Add cache entry
func (cache *Cache) Add(key, value string) {
	cache.data[key] = value
	cache.updated = true
}

// Get cache entry by key
func (cache *Cache) Get(key string) string {
	return cache.data[key]
}

// Write persists cache to disk
func (cache *Cache) Write() error {
	if !cache.updated || cache.fake {
		return nil
	}

	jsonString, err := json.Marshal(cache.data)
	if err != nil {
		return err
	}

	err = os.WriteFile(cache.path, jsonString, 0600)
	if err != nil {
		return err
	}

	return nil
}
