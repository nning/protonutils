package steam

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_GetGames(t *testing.T) {
	s, _ := New("", testConfig, false)

	games, err := s.LocalConfig.GetGames(false)
	assert.Empty(t, err)
	assert.Equal(t, 36, len(games))
}

func Test_GetGamesSteamOS(t *testing.T) {
	s, _ := New("", testConfig, false)

	games, err := s.LocalConfig.GetGames(true)
	assert.Empty(t, err)
	assert.Equal(t, 245, len(games))
}
