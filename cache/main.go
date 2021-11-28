package cache

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
	"time"
)

// Cache represents simple in-memory key/value store that can be persisted
type Cache struct {
	path    string
	data    map[string]Value
	updated bool
	maxAge  int64 // in seconds
}

// Value represents value in cache
type Value struct {
	Name      string `json:"name"`
	Valid     bool   `json:"valid"`
	UpdatedAt int64  `json:"updatedAt"` // in µs
}

// New instantiates new Cache
//   * name is the name of the cache file
//   * maxAge controls amount of seconds after which cache returns no entry even though an old one exists
func New(name string, maxAge int64) (*Cache, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	home := user.HomeDir
	dirPath := path.Join(home, ".cache", "protonutils")

	err = os.MkdirAll(dirPath, 0700)
	if err != err {
		return nil, err
	}

	filePath := path.Join(dirPath, name+".json")

	cache := &Cache{
		path:   filePath,
		data:   make(map[string]Value),
		maxAge: maxAge * 1000000, // in µs
	}

	if maxAge != 0 {
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
	cache.data[id] = Value{name, valid, time.Now().UnixMicro()}
	cache.updated = true
}

// Get cache entry by key
func (cache *Cache) Get(id string) (string, bool) {
	entry := cache.data[id]

	if cache.maxAge >= 0 && entry.UpdatedAt < time.Now().UnixMicro()-cache.maxAge {
		return "", false
	}

	return entry.Name, entry.Valid
}

// Write persists cache to disk
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

// Dump returns cache data
func (cache *Cache) Dump() map[string]Value {
	return cache.data
}
