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
	data    map[string]Value
	updated bool
	fake    bool
}

type Value struct {
	Name  string `json:"name"`
	Valid bool   `json:"valid"`
}

// New instantiates new Cache
func New(name string, fake bool) (*Cache, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	home := user.HomeDir
	p := path.Join(home, ".cache", name+".json")

	cache := &Cache{
		path: p,
		data: make(map[string]Value),
		fake: fake,
	}

	if !fake {
		f, err := os.ReadFile(cache.path)
		if err == nil {
			// Ignore errors; cache will be overwritten on Write
			json.Unmarshal(f, &cache.data)
		}
	}

	return cache, nil
}

// Add cache entry
func (cache *Cache) Add(id, name string, valid bool) {
	cache.data[id] = Value{name, valid}
	cache.updated = true
}

// Get cache entry by key
func (cache *Cache) Get(id string) (string, bool) {
	entry := cache.data[id]
	return entry.Name, entry.Valid
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

func (cache *Cache) Dump() map[string]Value {
	return cache.data
}
