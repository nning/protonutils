package steam

import (
	"encoding/json"
	"io"
	"net/http"

	. "github.com/nning/list_proton_versions"
)

type JsonResponse map[string]JsonAppData

type JsonAppData struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

func (self *Steam) GetName(id string) string {
	name := self.cache.Get(id)

	if name != "" {
		return name
	}

	res, err := http.Get("https://store.steampowered.com/api/appdetails/?appids=" + id)
	PanicOnError(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	PanicOnError(err)

	data := make(JsonResponse)
	err = json.Unmarshal(body, &data)
	PanicOnError(err)

	name = data[id].Data.Name
	val := name
	if val == "" {
		val = "💩"
	}

	self.cache.Add(id, val)

	return name
}

func (self *Steam) SaveCache() {
	self.cache.Write()
}
