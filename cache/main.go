package cache

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type Cache struct {
	path    string
	data    map[string]string
	updated bool
}

func New(name string) (*Cache, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	home := user.HomeDir
	p := path.Join(home, ".cache", name+".json")

	cache := &Cache{
		path: p,
		data: make(map[string]string),
	}

	f, err := os.ReadFile(cache.path)
	if err == nil {
		err = json.Unmarshal(f, &cache.data)
		if err != nil {
			return nil, err
		}
	}

	return cache, nil
}

func (cache *Cache) Add(key, value string) {
	cache.data[key] = value
	cache.updated = true
}

func (cache *Cache) Get(key string) string {
	return cache.data[key]
}

func (cache *Cache) Write() error {
	if !cache.updated {
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
