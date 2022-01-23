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
		n, err := s.AppInfo.GetName(id)
		assert.Empty(t, err)
		assert.Equal(t, name, n)
	}
}

func Test_GetNextEntryStart(t *testing.T) {
	s, err := New("", testSteamRoot, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	pos1 := 56
	pos2 := 143

	i := s.AppInfo.GetNextEntryStart(0)
	assert.Equal(t, pos1, i)

	n, err := ParseAppInfoEntry(s.AppInfo.Bytes[i:])
	assert.Empty(t, err)
	assert.NotEmpty(t, n)

	i = s.AppInfo.GetNextEntryStart(pos1 + 2)
	assert.Equal(t, pos2, i)

	n, err = ParseAppInfoEntry(s.AppInfo.Bytes[i:])
	assert.Empty(t, err)
	assert.NotEmpty(t, n)

	type v struct {
		Category DeckCompatibilityCategory
	}

	tests := map[string]*v{
		"1091500": {DeckCompatibilityUnknown},     // Cyberpunk 2077
		"1113000": {DeckCompatibilityUnsupported}, // Persona 4 Golden
		"1426210": {DeckCompatibilityPlayable},    // It Takes Two
		"292030":  {DeckCompatibilityPlayable},    // The Witcher 3
		"205100":  {DeckCompatibilityVerified},    // Dishonored
		"620":     {DeckCompatibilityVerified},    // Portal 2
		"1190460": {DeckCompatibilityVerified},    // Death Stranding
	}

	i = 0
	for {
		k := s.AppInfo.GetNextEntryStart(i)
		if k < 0 {
			break
		}

		n, err = ParseAppInfoEntry(s.AppInfo.Bytes[k:])
		assert.Empty(t, err)
		assert.NotEmpty(t, n)

		appID := n.FirstByName("appid").String()

		if tests[appID] != nil {
			// x, err := n.MarshalText()
			// assert.Empty(t, err)
			// os.WriteFile("debug-"+appID+".vdf", x, 0600)

			cn := n.FirstByName("common").FirstByName("steam_deck_compatibility")
			// compatibility, _ := cn.MarshalText()
			// fmt.Println("compatibility", string(compatibility))

			c := GetDeckCompatibility(cn)
			// fmt.Println(c)

			assert.Equal(t, tests[appID].Category, c.Category)
		}

		i = k + 2
	}
}
