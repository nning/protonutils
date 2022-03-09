package steam

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_GetViewedSteamPlay(t *testing.T) {
	s, _ := New("", testSteamRoot, false)

	games, err := s.LocalConfig.GetViewedSteamPlay()
	assert.Empty(t, err)
	assert.Equal(t, 36, len(games))
}
