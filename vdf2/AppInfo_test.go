package vdf2

import (
	"fmt"
	"os"
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

		if appID == "205100" { // Dishonored
			x, err := n.MarshalText()
			assert.Empty(t, err)
			os.WriteFile("debug.vdf", x, 0600)

			compatibility, _ := n.FirstByName("common").FirstByName("steam_deck_compatibility").MarshalText()
			fmt.Println("compatibility", string(compatibility))

			c := GetDeckCompatibility(n.FirstByName("common").FirstByName("steam_deck_compatibility"))
			fmt.Println(c)
		}

		i = k + 2
	}
}
