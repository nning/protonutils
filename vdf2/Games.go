package vdf2

import (
	"sort"

	"github.com/nning/protonutils/steam"
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

// GetGameData returns intialized Game struct by app ID
func GetGameData(s *steam.Steam, gameID string) (*Game, bool, error) {
	// TODO Implement GetGameData with vdf2
	g, isValid, err := s.GetGameData(gameID)
	if err != nil {
		return nil, isValid, err
	}

	return &Game{
		ID:          g.ID,
		Name:        g.Name,
		IsInstalled: g.IsInstalled,
		IsShortcut:  g.IsShortcut,
	}, isValid, nil
}
