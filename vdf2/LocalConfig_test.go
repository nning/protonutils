package vdf2

import (
	"github.com/nning/protonutils/steam"
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_GetViewedSteamPlay(t *testing.T) {
	s, err := steam.New("", testSteamRoot, true)
	assert.Empty(t, err)

	vdf, err := GetLocalConfig(s)
	assert.Empty(t, err)

	games, err := vdf.GetViewedSteamPlay()
	assert.Empty(t, err)
	assert.Equal(t, 36, len(games))
}
