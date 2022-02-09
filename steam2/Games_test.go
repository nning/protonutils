package steam2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetName(t *testing.T) {
	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	games := map[string]string{
		"292030":  "The Witcher 3: Wild Hunt",
		"377160":  "Fallout 4",
		"403640":  "Dishonored 2",
		"614570":  "Dishonored®: Death of the Outsider™ ",
		"826630":  "Iron Harvest",
		"1091500": "Cyberpunk 2077",
		"1151640": "Horizon Zero Dawn",
		"1174180": "Red Dead Redemption 2",
		"1328670": "Mass Effect™ Legendary Edition",
		"813780":  "Age of Empires II: Definitive Edition",
	}

	for id, name := range games {
		n, err := s.GetName(id)
		assert.Empty(t, err)
		assert.Equal(t, name, n)
	}
}

func Test_GetGame(t *testing.T) {
	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	game, err := s.GetGame("1091500")
	assert.Empty(t, err)
	assert.NotEmpty(t, game)

	assert.Equal(t, "Cyberpunk 2077", game.Name)
	assert.Equal(t, DeckCompatibilityUnknown, game.DeckCompatibility.Category)

	game, err = s.GetGame("205100")
	assert.Empty(t, err)
	assert.NotEmpty(t, game)

	assert.Equal(t, "Dishonored", game.Name)
	assert.Equal(t, DeckCompatibilityVerified, game.DeckCompatibility.Category)
}

func Test_GetShortcutName(t *testing.T) {
	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.Shortcuts)

	n, err := s.GetShortcutName("3228583970")
	assert.Empty(t, err)
	assert.Equal(t, "Kena - Bridge of Spirits", n)

	n, err = s.GetShortcutName("2977655160")
	assert.Empty(t, err)
	assert.Equal(t, "", n)
}

func Test_GetGameData(t *testing.T) {
	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	games := map[string]string{
		"292030":     "The Witcher 3: Wild Hunt",
		"377160":     "Fallout 4",
		"403640":     "Dishonored 2",
		"614570":     "Dishonored®: Death of the Outsider™ ",
		"826630":     "Iron Harvest",
		"1091500":    "Cyberpunk 2077",
		"1151640":    "Horizon Zero Dawn",
		"1174180":    "Red Dead Redemption 2",
		"1328670":    "Mass Effect™ Legendary Edition",
		"813780":     "Age of Empires II: Definitive Edition",
		"3228583970": "Kena - Bridge of Spirits",
	}

	for id, name := range games {
		g, isValid, err := s.GetGameData(id)
		assert.Empty(t, err)
		assert.True(t, isValid)
		assert.NotEmpty(t, g)
		assert.Equal(t, name, g.Name)
	}
}

func Test_GetGameData_MissingShortcut(t *testing.T) {
	s, err := New("", testSteamRoot, true)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	g, isValid, err := s.GetGameData("2977655160")
	assert.Empty(t, err)
	assert.False(t, isValid)

	assert.NotEmpty(t, g)
	assert.Equal(t, "", g.Name)
	assert.False(t, g.IsInstalled)
	assert.True(t, g.IsShortcut)
}
