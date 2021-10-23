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

const InvalidId = "💩"

func (s *Steam) GetName(id string) string {
	name := s.cache.Get(id)

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
		val = InvalidId
	}

	s.cache.Add(id, val)

	return name
}

func (s *Steam) GetGameData(id string) *GameData {
	return &GameData{id, s.IsInstalled(id)}
}

func (s *Steam) SaveCache() {
	s.cache.Write()
}
