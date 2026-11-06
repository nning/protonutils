package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/nning/protonutils/utils"
)

// Cache represents simple in-memory key/value store that can be persisted
type Cache struct {
	path    string
	data    map[string]Value
	updated bool
	maxAge  time.Duration
}

// Value represents value in cache
type Value struct {
	Name      string `json:"name"`
	Valid     bool   `json:"valid"`
	UpdatedAt int64  `json:"updatedAt"` // in µs
}

// Age returns time elapsed since the value has been updated
func (v Value) Age() time.Duration {
	return time.Since(time.UnixMicro(v.UpdatedAt))
}

// New instantiates new Cache
//   - name is the name of the cache file
//   - maxAge controls the duration after which cache returns no entry even
//     though an old one exists (negative values disable aging, 0 ignores the
//     persisted cache and only returns entries added during the current run)
func New(name string, maxAge time.Duration) (*Cache, error) {
	dir := utils.GetCacheDir()

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return nil, err
	}

	cache := &Cache{
		path:   filepath.Join(dir, name+".json"),
		data:   make(map[string]Value),
		maxAge: maxAge,
	}

	if maxAge != 0 {
		f, err := os.ReadFile(cache.path)
		if err == nil {
			// Ignore errors; cache will be overwritten on Write
			if json.Unmarshal(f, &cache.data) != nil {
				cache.data = make(map[string]Value)
			}
		}
	}

	return cache, nil
}

// Add cache entry
func (cache *Cache) Add(id, name string, valid bool) {
	cache.data[id] = Value{name, valid, time.Now().UnixMicro()}
	cache.updated = true
}

// Lookup returns the cache entry for key id and whether a non-expired entry
// exists
func (cache *Cache) Lookup(id string) (Value, bool) {
	entry, ok := cache.data[id]
	if !ok {
		return Value{}, false
	}

	if cache.maxAge > 0 && entry.Age() > cache.maxAge {
		return Value{}, false
	}

	return entry, true
}

// Get returns name and validity of cache entry by key (empty string and false
// if there is no non-expired entry)
func (cache *Cache) Get(id string) (string, bool) {
	entry, ok := cache.Lookup(id)
	if !ok {
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

	tmp := cache.path + ".tmp"
	err = os.WriteFile(tmp, jsonString, 0600)
	if err != nil {
		return err
	}

	err = os.Rename(tmp, cache.path)
	if err != nil {
		return err
	}

	cache.updated = false
	return nil
}

// Dump returns cache data
func (cache *Cache) Dump() map[string]Value {
	return cache.data
}
