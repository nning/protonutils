package steam

import (
	"sort"
)

// Game represents Steam game or shortcut
type Game struct {
	ID          string `json:"appID"`
	Name        string `json:"-"`
	IsInstalled bool   `json:"isInstalled"`
	IsShortcut  bool   `json:"isShortcut"`
}

// Games maps game name to Game (app ID, install status)
type Games map[string]*Game

func (s *Steam) addGame(versionID, versionName, gameID string, isDefault bool) (*Game, error) {
	game, valid, err := s.GetGameData(gameID)
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

	s.CompatToolVersions[versionName].Games[game.Name] = game

	return game, nil
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
