package steam

import (
	"sort"
)

type gameData struct {
	ID          string `json:"appID"`
	IsInstalled bool   `json:"isInstalled"`
	IsShortcut  bool   `json:"isShortcut"`
}

// Games maps game name to gameData (app ID, install status)
type Games map[string]*gameData

func (s *Steam) addGame(versionID, versionName, gameID string, isDefault bool) (*gameData, error) {
	name, data, valid, err := s.getNameAndGameData(gameID)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, nil
	}

	if s.CompatToolVersions[versionName] == nil {
		s.CompatToolVersions[versionName] = &Version{
			ID:        versionID,
			Name:      versionName,
			Games:     make(Games),
			IsDefault: isDefault,
		}
	}

	if s.CompatToolVersions[versionName].ID == "" {
		s.CompatToolVersions[versionName].IsDefault = true
	}

	s.CompatToolVersions[versionName].Games[name] = data

	return data, nil
}

func (games Games) includesID(id string) bool {
	for _, data := range games {
		if data.ID == id {
			return true
		}
	}

	return false
}

// Sort returns slice of alphabetically sorted Game names
func (games Games) Sort() []string {
	keys := make([]string, len(games))

	i := 0
	for key := range games {
		keys[i] = key
		i++
	}

	sort.Strings(keys)
	return keys
}
