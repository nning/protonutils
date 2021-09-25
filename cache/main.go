package cache

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path"

	. "github.com/nning/list_proton_versions"
)

type Cache struct {
	path    string
	data    map[string]string
	updated bool
}

func panicOnError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func New(name string) *Cache {
	user, err := user.Current()
	PanicOnError(err)

	home := user.HomeDir
	p := path.Join(home, ".cache", name+".json")

	cache := &Cache{
		path: p,
		data: make(map[string]string),
	}

	f, err := os.ReadFile(cache.path)
	if err == nil {
		err = json.Unmarshal(f, &cache.data)
		PanicOnError(err)
	}

	return cache
}

func (cache *Cache) Add(key, value string) {
	cache.data[key] = value
	cache.updated = true
}

func (cache *Cache) Get(key string) string {
	return cache.data[key]
}

func (cache *Cache) Write() {
	if !cache.updated {
		return
	}

	jsonString, err := json.Marshal(cache.data)
	PanicOnError(err)

	err = os.WriteFile(cache.path, jsonString, 0600)
	PanicOnError(err)
}
