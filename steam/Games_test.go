package steam

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

func Test_GetAppIDAndName(t *testing.T) {
	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s)

	s.ReadCompatTools()

	// Exact name match but lowercase
	results := s.GetAppIDAndNames("disco elysium")
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "632470", results[0][0])
	assert.Equal(t, "Disco Elysium", results[0][1])

	// ID match
	results = s.GetAppIDAndNames("632470")
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "632470", results[0][0])
	assert.Equal(t, "Disco Elysium", results[0][1])

	// Prefix and lowercase match
	results = s.GetAppIDAndNames("cyberpunk")
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "1091500", results[0][0])
	assert.Equal(t, "Cyberpunk 2077", results[0][1])

	// Several matches
	results = s.GetAppIDAndNames("fallout")
	assert.Equal(t, 4, len(results))

	expected := [][]string{
		{"22370", "Fallout 3 - Game of the Year Edition"},
		{"22380", "Fallout: New Vegas"},
		{"377160", "Fallout 4"},
		{"1151340", "Fallout 76"},
	}

	for _, res := range results {
		var x []string

		for _, ex := range expected {
			if res[0] == ex[0] {
				x = res
			}
		}

		assert.Equal(t, 2, len(x))
		assert.Equal(t, x[0], res[0])
		assert.Equal(t, x[1], res[1])
	}
}

func Test_GetGameVersion(t *testing.T) {
	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s)

	err = s.ReadCompatTools()
	assert.Empty(t, err)

	tool := s.GetGameVersion("1252330")
	assert.NotEmpty(t, tool)
	assert.Equal(t, "Proton-7.0rc6-GE-1", tool.ID)
	assert.Equal(t, "Proton-7.0rc6-GE-1", tool.Name)

	tool = s.GetGameVersion("1222140")
	assert.NotEmpty(t, tool)
	assert.Equal(t, "", tool.ID)
	assert.Equal(t, "Proton 6.3-8 (Default)", tool.Name)

	tool = s.GetGameVersion("11111111")
	assert.Empty(t, tool)

	tool = s.GetGameVersion("")
	assert.Empty(t, tool)
}
