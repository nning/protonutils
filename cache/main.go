package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

type Cache struct {
	Path string
	Data map[string]string
}

func panicOnError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func New(name string) *Cache {
	user, err := user.Current()
	panicOnError(err)

	home := user.HomeDir
	p := path.Join(home, ".cache", name+".json")

	cache := &Cache{
		Path: p,
		Data: make(map[string]string),
	}

	f, err := os.ReadFile(cache.Path)
	if err == nil {
		err = json.Unmarshal(f, &cache.Data)
		panicOnError(err)
	} else {
		fmt.Println("Create file on write")
	}

	return cache
}

func (cache *Cache) Add(key, value string) {
	cache.Data[key] = value
}

func (cache *Cache) Get(key string) string {
	return cache.Data[key]
}

func (cache *Cache) Write() {
	jsonString, err := json.Marshal(cache.Data)
	panicOnError(err)

	err = os.WriteFile(cache.Path, jsonString, 0600)
	panicOnError(err)
}