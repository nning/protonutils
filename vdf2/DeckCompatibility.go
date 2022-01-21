package vdf2

import (
	"github.com/BenLubar/vdf"
)

type DeckCompatibilityCategory int

const (
	Unknown DeckCompatibilityCategory = iota
	Unsupported
	Playable
	Verified
)

type DeckCompatibility struct {
	Category DeckCompatibilityCategory
}

func GetDeckCompatibility(n *vdf.Node) *DeckCompatibility {
	category := DeckCompatibilityCategory(n.FirstByName("category").Int())
	c := &DeckCompatibility{Category: category}

	return c
}
