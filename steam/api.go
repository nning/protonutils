package steam

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type jsonAppData struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

// GetNameFromAPI returns name of game by app id using the Steam API
func (s *Steam) GetNameFromAPI(id string) (string, error) {
	log.Debug("Steam.GetNameFromAPI(", id, ")")

	res, err := http.Get("https://store.steampowered.com/api/appdetails/?appids=" + id)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	data := make(map[string]jsonAppData)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data[id].Data.Name, nil
}
