package steam

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// StoreAPIURL is the base URL of the Steam store API (variable for tests)
var StoreAPIURL = "https://store.steampowered.com/api/appdetails/"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type jsonAppData struct {
	Success bool `json:"success"`
	Data    struct {
		Name string `json:"name"`
	} `json:"data"`
}

// GetNameFromAPI returns name of game by app id using the Steam store API
func (s *Steam) GetNameFromAPI(id string) (string, error) {
	slog.Debug("Steam.GetNameFromAPI", "id", id)

	res, err := httpClient.Get(StoreAPIURL + "?appids=" + url.QueryEscape(id))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("steam store API: %v", res.Status)
	}

	data := make(map[string]jsonAppData)
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data[id].Data.Name, nil
}
