package steam

import (
	"sort"
)

type GameData struct {
	Id          string `json:"id"`
	IsInstalled bool   `json:"isInstalled"`
}

type Games map[string]*GameData

func (s *Steam) AddGame(version, id string) (*GameData, error) {
	name, err := s.GetName(id)
	if err != nil {
		return nil, err
	}

	if name != InvalidId {
		if s.CompatToolVersions[version] == nil {
			s.CompatToolVersions[version] = make(Games)
		}

		data, err := s.GetGameData(id)
		if err != nil {
			return nil, err
		}

		s.CompatToolVersions[version][name] = data
		return data, nil
	}

	return nil, nil
}

func (games Games) IncludesId(id string) bool {
	for _, data := range games {
		if data.Id == id {
			return true
		}
	}

	return false
}

func (games Games) Sort() []string {
	var keys []string

	for key := range games {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}
