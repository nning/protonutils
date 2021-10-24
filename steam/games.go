package steam

import (
	"sort"
)

type GameData struct {
	Id          string `json:"id"`
	IsInstalled bool   `json:"isInstalled"`
}

type Games map[string]*GameData

func (s *Steam) AddGame(version, id string) {
	name := s.GetName(id)

	if name != InvalidId {
		if s.CompatToolVersions[version] == nil {
			s.CompatToolVersions[version] = make(Games)
		}

		s.CompatToolVersions[version][name] = s.GetGameData(id)
	}
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
