package steam

import (
	"sort"
)

type gameData struct {
	ID          string `json:"appID"`
	IsInstalled bool   `json:"isInstalled"`
}

type games map[string]*gameData

func (s *Steam) addGame(version, id string) (*gameData, error) {
	name, valid, err := s.getName(id)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, nil
	}

	if s.CompatToolVersions[version] == nil {
		s.CompatToolVersions[version] = make(games)
	}

	data, err := s.getGameData(id)
	if err != nil {
		return nil, err
	}

	s.CompatToolVersions[version][name] = data
	return data, nil
}

func (games games) includesID(id string) bool {
	for _, data := range games {
		if data.ID == id {
			return true
		}
	}

	return false
}

// Sort returns slice of alphabetically sorted Game names
func (games games) Sort() []string {
	keys := make([]string, len(games))

	i := 0
	for key := range games {
		keys[i] = key
		i++
	}

	sort.Strings(keys)
	return keys
}
