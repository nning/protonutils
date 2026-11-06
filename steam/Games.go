package steam

import (
	"encoding/json"
	"errors"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/nning/protonutils/cache"
)

// negativeCacheMaxAge is the time after which app IDs without known name are
// looked up again
const negativeCacheMaxAge = 7 * 24 * time.Hour

// Game represents Steam game or shortcut
type Game struct {
	ID                string            `json:"appID"`
	Name              string            `json:"name"`
	IsInstalled       bool              `json:"isInstalled"`
	IsShortcut        bool              `json:"isShortcut"`
	DeckCompatibility DeckCompatibility `json:"deckCompatibility"`

	// Paths and sizes (in bytes) are only set on request (see
	// Steam.GetGamePaths)
	CompatdataPath string `json:"compatdataPath,omitempty"`
	CompatdataSize uint64 `json:"compatdataSize,omitempty"`
	InstallPath    string `json:"installPath,omitempty"`
	InstallSize    uint64 `json:"installSize,omitempty"`
}

// Games maps app ID to Game
type Games map[string]*Game

// AmbiguousNameError is returned if game search by app ID or name leads to
// several results
type AmbiguousNameError struct {
	Matches []AppIDAndName
}

// Error returns description string
func (err *AmbiguousNameError) Error() string {
	return "Ambiguous name, try using app ID"
}

// AppIDAndName holds app ID and name of a game search result
type AppIDAndName struct {
	ID   string
	Name string
}

// MarshalJSON serializes games as object with game names as keys (duplicate
// names are suffixed with the app ID)
func (games Games) MarshalJSON() ([]byte, error) {
	count := make(map[string]int)
	for _, game := range games {
		count[game.Name]++
	}

	byName := make(map[string]*Game, len(games))
	for _, game := range games {
		key := game.Name
		if count[key] > 1 || key == "" {
			key += " (" + game.ID + ")"
		}

		byName[key] = game
	}

	return json.Marshal(byName)
}

// Sort returns slice of app IDs sorted by game name (and app ID)
func (games Games) Sort() []string {
	ids := make([]string, 0, len(games))
	for id := range games {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		a, b := games[ids[i]], games[ids[j]]
		if a.Name != b.Name {
			return a.Name < b.Name
		}

		return a.ID < b.ID
	})

	return ids
}

// SortByDeckCompatibilityCategory returns games with deck compatibility info
// sorted by category
func (games Games) SortByDeckCompatibilityCategory() []Game {
	var sorted []Game

	for _, id := range games.Sort() {
		game := games[id]
		if game.DeckCompatibility.Category.String() != "" {
			sorted = append(sorted, *game)
		}
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].DeckCompatibility.Category < sorted[j].DeckCompatibility.Category
	})

	return sorted
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

// Includes returns whether app id is included in games
func (games Games) Includes(id string) bool {
	_, ok := games[id]
	return ok
}

// GetGameData returns initialized Game struct by app ID and whether a name was
// found
func (s *Steam) GetGameData(id string) (*Game, bool, error) {
	isShortcut := IsShortcut(id)
	isInstalled := s.LibraryConfig.IsInstalled(id) || isShortcut

	game := &Game{
		ID:          id,
		IsInstalled: isInstalled,
		IsShortcut:  isShortcut,
	}

	if isShortcut {
		game.Name = s.Shortcuts.GetName(id)
		if game.Name == "" {
			game.IsInstalled = false
		}

		// Cache name for searching by name (GetAppIDAndNames)
		valid := game.Name != ""
		s.AppidCache.Add(id, game.Name, valid)

		return game, valid, nil
	}

	if entry, ok := s.AppidCache.Lookup(id); ok {
		if entry.Valid && entry.Name != "" {
			game.Name = entry.Name
			slog.Debug("Steam.GetGameData (cache hit)", "id", id, "name", entry.Name)
			return game, true, nil
		}

		if !entry.Valid && entry.Age() < negativeCacheMaxAge {
			slog.Debug("Steam.GetGameData (negative cache hit)", "id", id)
			return game, false, nil
		}
	}

	g, err := s.GetGame(id)
	if err != nil {
		// A single unparsable appinfo.vdf entry must not break everything
		slog.Debug("Steam.GetGameData: ignoring appinfo.vdf entry", "id", id, "err", err)
		g = nil
	}

	if g != nil && g.Name != "" {
		game = g
	} else if id != "0" {
		if s.apiUnavailable {
			return game, false, nil
		}

		game.Name, err = s.GetNameFromAPI(id)
		if err != nil {
			// Do not fail (or cache) if the store API is not reachable
			slog.Debug("Steam.GetGameData: store API unavailable", "id", id, "err", err)
			s.apiUnavailable = true
			return game, false, nil
		}
	}

	slog.Debug("Steam.GetGameData (cache miss)", "id", id, "name", game.Name)

	valid := game.Name != ""
	s.AppidCache.Add(id, game.Name, valid)

	return game, valid, nil
}

// GetGame returns Game struct by app id (nil if app ID is not in appinfo.vdf)
func (s *Steam) GetGame(id string) (*Game, error) {
	slog.Debug("Steam.GetGame", "id", id)

	n, err := s.AppInfo.Get(id)
	if n == nil || err != nil {
		return nil, err
	}

	common := n.FirstByName("common")

	return &Game{
		ID:                id,
		Name:              common.FirstByName("name").String(),
		IsInstalled:       s.LibraryConfig.IsInstalled(id),
		IsShortcut:        IsShortcut(id),
		DeckCompatibility: *GetDeckCompatibility(common.FirstByName("steam_deck_compatibility")),
	}, nil
}

// GetName returns name of game by app id
func (s *Steam) GetName(id string) (string, error) {
	game, err := s.GetGame(id)
	if game == nil || err != nil {
		return "", err
	}

	return game.Name, nil
}

// GetShortcutName returns name of shortcut by app id
func (s *Steam) GetShortcutName(id string) (string, error) {
	return s.Shortcuts.GetName(id), nil
}

// GetAppIDAndNames returns app IDs and proper names for app ID or name search
// (from AppidCache cache), sorted by name
func (s *Steam) GetAppIDAndNames(idOrName string) []AppIDAndName {
	return SearchAppIDAndNames(s.AppidCache, idOrName)
}

// SearchAppIDAndNames returns app IDs and proper names for app ID or name
// search (case-insensitive name prefix) in an app ID cache, sorted by name
func SearchAppIDAndNames(c *cache.Cache, idOrName string) []AppIDAndName {
	results := make([]AppIDAndName, 0)
	b := strings.ToLower(strings.TrimSpace(idOrName))

	for id, value := range c.Dump() {
		if !value.Valid || value.Name == "" {
			continue
		}

		a := strings.ToLower(value.Name)

		if strings.HasPrefix(a, b) || id == idOrName {
			results = append(results, AppIDAndName{id, value.Name})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Name != results[j].Name {
			return results[i].Name < results[j].Name
		}

		return results[i].ID < results[j].ID
	})

	return results
}

// GetAppIDAndName returns app ID and proper name for app ID or name search
// (from AppidCache cache). If several games match, an exact (case-insensitive)
// match of app ID or name is preferred.
func (s *Steam) GetAppIDAndName(idOrName string) (string, string, error) {
	results := s.GetAppIDAndNames(idOrName)

	if len(results) == 0 {
		return "", "", errors.New("App ID or name not found")
	}

	if len(results) == 1 {
		return results[0].ID, results[0].Name, nil
	}

	var exact []AppIDAndName
	for _, r := range results {
		if r.ID == idOrName || strings.EqualFold(r.Name, strings.TrimSpace(idOrName)) {
			exact = append(exact, r)
		}
	}

	if len(exact) == 1 {
		return exact[0].ID, exact[0].Name, nil
	}

	return "", "", &AmbiguousNameError{Matches: results}
}

// GetGameVersion returns compatibility tool used by game with app ID id
func (s *Steam) GetGameVersion(id string) *CompatTool {
	for _, tool := range s.CompatTools {
		if tool.Games.Includes(id) {
			return tool
		}
	}

	return nil
}
