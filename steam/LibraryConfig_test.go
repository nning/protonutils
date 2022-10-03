package steam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LibraryConfig(t *testing.T) {
	s, _ := New("", testConfig, false)

	p := s.LibraryConfig.GetLibraryPathByID("1091500")
	assert.Equal(t, "/mnt/games/Shared Steam Library", p)

	x := s.LibraryConfig.IsInstalled("1091500")
	assert.True(t, x)
}
