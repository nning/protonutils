package steam2

import (
	"encoding/json"
	"io"
	"net/http"
)

type jsonAppData struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

func (s *Steam) GetNameFromAPI(id string) (string, error) {
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
