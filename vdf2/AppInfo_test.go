package vdf2

import (
	"testing"

	"github.com/nning/protonutils/steam"
	"github.com/stretchr/testify/assert"
)

func Test_GetNextEntryStart(t *testing.T) {
	s, err := steam.New("", testSteamRoot, true)
	assert.Empty(t, err)

	ai, err := GetAppInfo(s)
	assert.Empty(t, err)
	assert.NotEmpty(t, ai)

	pos1 := 56
	pos2 := 143

	i := ai.GetNextEntryStart(0)
	assert.Equal(t, pos1, i)

	n, err := ParseAppInfoEntry(ai.Bytes[i:])
	assert.Empty(t, err)
	assert.NotEmpty(t, n)

	i = ai.GetNextEntryStart(pos1 + 2)
	assert.Equal(t, pos2, i)

	n, err = ParseAppInfoEntry(ai.Bytes[i:])
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
		k := ai.GetNextEntryStart(i)
		if k < 0 {
			break
		}

		n, err = ParseAppInfoEntry(ai.Bytes[k:])
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
