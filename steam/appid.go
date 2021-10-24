package steam

import (
	"encoding/json"
	"io"
	"net/http"
)

type JsonResponse map[string]JsonAppData

type JsonAppData struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

const InvalidId = "💩"

func (s *Steam) GetName(id string) (string, error) {
	name := s.appidCache.Get(id)

	if name != "" {
		return name, nil
	}

	res, err := http.Get("https://store.steampowered.com/api/appdetails/?appids=" + id)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	data := make(JsonResponse)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	name = data[id].Data.Name
	val := name
	if val == "" {
		val = InvalidId
	}

	s.appidCache.Add(id, val)

	return name, nil
}

func (s *Steam) GetGameData(id string) (*GameData, error) {
	isInstalled, err := s.IsInstalled(id)
	if err != nil {
		return nil, err
	}

	return &GameData{id, isInstalled}, nil
}

func (s *Steam) SaveCache() error {
	return s.appidCache.Write()
}
