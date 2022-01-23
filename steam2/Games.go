package steam2

import (
	"sort"
)

// Game represents Steam game or shortcut
type Game struct {
	ID                string            `json:"appID"`
	Name              string            `json:"-"`
	IsInstalled       bool              `json:"isInstalled"`
	IsShortcut        bool              `json:"isShortcut"`
	DeckCompatibility DeckCompatibility `json:"-"`
}

// Games maps game name to Game (app ID, install status)
type Games map[string]*Game

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

// CountInstalled returns count of installed games
func (games Games) CountInstalled() int {
	i := 0

	for _, game := range games {
		if game.IsInstalled {
			i++
		}
	}

	return i
}

// Includes returns whether appID is included in g
func (games Games) Includes(appID string) bool {
	for _, data := range games {
		if data.ID == appID {
			return true
		}
	}

	return false
}

// GetGameData returns initialized Game struct by app ID
func (s *Steam) GetGameData(id string) (*Game, bool, error) {
	isShortcut := IsShortcut(id)
	isInstalled := s.LibraryConfig.IsInstalled(id)

	name, valid := s.AppidCache.Get(id)
	if name != "" && valid {
		data := &Game{
			ID:          id,
			Name:        name,
			IsInstalled: isInstalled,
			IsShortcut:  isShortcut,
		}

		return data, true, nil
	}

	var err error

	if isShortcut {
		// TODO Get name from shortcuts
		name = ""
	} else {
		name, err = s.AppInfo.GetName(id)
		if err != nil {
			return nil, false, err
		}

		if name == "" {
			name, err = s.GetNameFromAPI(id)
			if err != nil {
				return nil, false, err
			}
		}
	}

	valid = name != ""
	s.AppidCache.Add(id, name, valid)

	return &Game{
		ID:          id,
		Name:        name,
		IsInstalled: isInstalled,
		IsShortcut:  isShortcut,
	}, valid, nil
}
