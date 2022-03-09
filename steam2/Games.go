package steam2

import (
	"errors"
	"sort"
	"strings"
)

// Game represents Steam game or shortcut
type Game struct {
	ID                string            `json:"appID"`
	Name              string            `json:"-"`
	IsInstalled       bool              `json:"isInstalled"`
	IsShortcut        bool              `json:"isShortcut"`
	DeckCompatibility DeckCompatibility `json:"deckCompatibility"`
}

// Games maps game name to Game (app ID, install status)
type Games map[string]*Game

// AmbiguousNameError is returned if game search by app ID or name leads to
// several results
type AmbiguousNameError struct{}

// Error returns description string
func (err *AmbiguousNameError) Error() string {
	return "Ambiguous name, try using app ID"
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

// Includes returns whether app id is included in g
func (games Games) Includes(id string) bool {
	for _, data := range games {
		if data.ID == id {
			return true
		}
	}

	return false
}

// GetGameData returns initialized Game struct by app ID
// TODO Rename/rewrite GetGameCached? How to handle non-cached data?
func (s *Steam) GetGameData(id string) (*Game, bool, error) {
	isShortcut := IsShortcut(id)
	isInstalled := s.LibraryConfig.IsInstalled(id) || isShortcut

	game := &Game{
		ID:          id,
		IsInstalled: isInstalled,
		IsShortcut:  isShortcut,
	}

	name, valid := s.AppidCache.Get(id)
	if name != "" && valid {
		game.Name = name
		return game, true, nil
	}

	var err error

	if isShortcut {
		game.Name, err = s.GetShortcutName(id)
		if err != nil {
			return nil, false, err
		}

		if game.Name == "" {
			game.IsInstalled = false
		}
	} else {
		game1, err := s.GetGame(id)
		if err != nil {
			return nil, false, err
		}

		if game1 != nil {
			game = game1
		} else if id != "0" {
			game.Name, err = s.GetNameFromAPI(id)
			if err != nil {
				return nil, false, err
			}
		}
	}

	valid = game.Name != ""
	s.AppidCache.Add(id, game.Name, valid)

	return game, valid, nil
}

// GetGame returns Game struct by app id
// If nameOnly is true, Game struct will only contain name
func (s *Steam) GetGame(id string, nameOnly ...bool) (*Game, error) {
	i, err := s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, id)
	if i < 0 || err != nil {
		return nil, err
	}

	n, err := ParseBinaryVdf(s.AppInfo.Bytes[i:])
	if err != nil && err.Error() != "vdf: unknown pack type 114" {
		return nil, err
	}

	name := n.FirstByName("common").FirstByName("name").String()
	game := &Game{
		ID:   id,
		Name: name,
	}

	if len(nameOnly) == 0 || (len(nameOnly) > 0 && nameOnly[0] == false) {
		game.IsInstalled = s.LibraryConfig.IsInstalled(id)
		game.IsShortcut = IsShortcut(id)

		cn := n.FirstByName("common").FirstByName("steam_deck_compatibility")
		game.DeckCompatibility = *GetDeckCompatibility(cn)
	}

	return game, nil
}

// GetName returns name of game by app id
func (s *Steam) GetName(id string) (string, error) {
	game, err := s.GetGame(id, true)
	if game == nil || err != nil {
		return "", err
	}

	return game.Name, err
}

// GetShortcutName returns name of shortcut by app id
func (s *Steam) GetShortcutName(id string) (string, error) {
	i, err := s.Shortcuts.GetNextEntryStartByID(0, InnerOffsetShortcuts, id)
	if i < 0 || err != nil {
		return "", nil
	}

	n, err := ParseBinaryVdf(s.Shortcuts.Bytes[i:])
	if err != nil {
		return "", err
	}

	return n.NextByName("AppName").String(), nil
}

// GetAppIDAndNames returns app IDs and proper names for name search (from
// AppidCache cache)
func (s *Steam) GetAppIDAndNames(idOrName string) [][]string {
	data := s.AppidCache.Dump()

	results := make([][]string, 0)

	for id, value := range data {
		a := strings.ToLower(value.Name)
		b := strings.ToLower(idOrName)

		if a == b || strings.HasPrefix(a, b) || id == idOrName {
			results = append(results, []string{id, value.Name})
		}
	}

	return results
}

func (s *Steam) GetAppIDAndName(idOrName string) (string, string, error) {
	idAndNames := s.GetAppIDAndNames(idOrName)

	l := len(idAndNames)
	if l == 0 {
		return "", "", errors.New("App ID or name not found")
	} else if l > 1 {
		return "", "", &AmbiguousNameError{}
	}

	return idAndNames[0][0], idAndNames[0][1], nil
}
