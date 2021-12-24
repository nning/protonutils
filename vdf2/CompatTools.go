package vdf2

import (
	"github.com/nning/protonutils/steam"
)

// CompatTools maps compatibility tool version IDs to CompatTool objects
// containing info like name and games.
type CompatTools map[string]*CompatTool

// CompatTool holds info about a compatibility tool (like human-readable name
// and a list of the games that are using it)
type CompatTool = steam.Version

// IsValid checks whether a version ID (v) exists in the CompatTools config
func (c CompatTools) IsValid(v string) bool {
	for id := range c {
		if id == v {
			return true
		}
	}

	return false
}

// Add adds an entry to the CompatTools config (by version id and name)
func (c CompatTools) Add(id, name string) {
	if c[id] != nil {
		return
	}

	c[id] = &CompatTool{
		ID:        id,
		Name:      name,
		IsDefault: false, // TODO
		Games:     make(steam.Games),
	}
}

// AddGame adds a Game entry to an existing CompatTool entry
func (c CompatTools) AddGame(id string, game *steam.Game) bool {
	if c[id] == nil {
		return false
	}

	c[id].Games[game.Name] = game
	return true
}
