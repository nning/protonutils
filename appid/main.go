package appid

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nning/list_proton_versions/cache"
)

type AppId struct {
	cache *cache.Cache
}

type JsonAppData struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

func panicOnError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func New() *AppId {
	c := cache.New("steam-appids")
	return &AppId{
		cache: c,
	}
}

func (self *AppId) GetName(id string) string {
	name := self.cache.Get(id)

	if name == "" {
		res, err := http.Get("https://store.steampowered.com/api/appdetails/?appids=" + id)
		panicOnError(err)

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		panicOnError(err)

		data := make(map[string]JsonAppData)
		err = json.Unmarshal(body, &data)
		panicOnError(err)

		name = data[id].Data.Name
		val := name
		if val == "" {
			val = "💩"
		}

		self.cache.Add(id, val)
	}

	return name
}

func (self *AppId) Write() {
	self.cache.Write()
}
