package steam

import (
	"github.com/BenLubar/vdf"
)

// DeckCompatibility holds deck compatibility info for a game
type DeckCompatibility struct {
	Category      DeckCompatibilityCategory `json:"category"`
	Tests         []string                  `json:"tests"`
	Configuration map[string]string         `json:"configuration"`
}

// String returns string that represents deck compatibility status
func (c DeckCompatibility) String() string {
	return c.Category.String()
}

// GetDeckCompatibility reads deck compatibility data from VDF node n
func GetDeckCompatibility(n *vdf.Node) *DeckCompatibility {
	category := DeckCompatibilityCategory(n.FirstByName("category").Int())
	tests := make([]string, 0) // TODO

	c := &DeckCompatibility{
		Category: category,
		Tests:    tests,
	}

	return c
}
