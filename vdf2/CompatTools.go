package vdf2

import (
	"github.com/nning/protonutils/steam"
)

type CompatTools map[string]*CompatTool

type CompatTool = steam.Version

func (c CompatTools) IsValid(v string) bool {
	for id := range c {
		if id == v {
			return true
		}
	}

	return false
}

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

func (c CompatTools) AddGame(id string, game *steam.Game) {
	c[id].Games[game.Name] = game
}
