package vdf2

import (
	"github.com/BenLubar/vdf"
)

type DeckCompatibility struct {
	Category DeckCompatibilityCategory
	Tests    []string
}

func (c DeckCompatibility) String() string {
	return "DeckCompatibility{" + c.Category.String() + "}"
}

func GetDeckCompatibility(n *vdf.Node) *DeckCompatibility {
	category := DeckCompatibilityCategory(n.FirstByName("category").Int())
	tests := make([]string, 0)

	c := &DeckCompatibility{
		Category: category,
		Tests:    tests,
	}

	return c
}
