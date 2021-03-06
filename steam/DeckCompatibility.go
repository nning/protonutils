package steam

import (
	"github.com/BenLubar/vdf"
)

// DeckCompatibility holds deck compatibility info for a game
type DeckCompatibility struct {
	Category DeckCompatibilityCategory `json:"category"`
	Tests    []string                  `json:"tests"`
}

// String returns string that represents deck compatibility status
func (c DeckCompatibility) String() string {
	return "DeckCompatibility{" + c.Category.String() + "}"
}

// GetDeckCompatibility reads deck compatibility data from VDF node n
func GetDeckCompatibility(n *vdf.Node) *DeckCompatibility {
	category := DeckCompatibilityCategory(n.FirstByName("category").Int())
	tests := make([]string, 0)

	c := &DeckCompatibility{
		Category: category,
		Tests:    tests,
	}

	return c
}
