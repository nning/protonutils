package steam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DeckCompatibility(t *testing.T) {
	s, _ := New("", testSteamRoot, true)

	games := map[string]string{
		"1252330": "Verified",    // DEATHLOOP
		"812140":  "Playable",    // Assassin's Creed Odyssey
		"813780":  "Unsupported", // Age of Empires II: Definitive Edition
		"1222140": "Unknown",     // Detroit: Become Human
	}

	for gameID, expected := range games {
		game, err := s.GetGame(gameID)
		assert.Empty(t, err)
		assert.Equal(t, expected, game.DeckCompatibility.Category.String())
	}
}
