package steam

import (
	"errors"
	"strconv"
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

	name, err := s.FindNameInAppInfo(id)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return "", err
	}

	if name == "" {
		name = InvalidId
	}

	s.appidCache.Add(id, name)
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
