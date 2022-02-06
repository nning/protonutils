package steam2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetCompatToolName(t *testing.T) {
	s, _ := New("", testSteamRoot, true)

	name, err := s.GetCompatToolName("proton_63")

	assert.Empty(t, err)
	assert.Equal(t, "Proton 6.3-8", name)

	name, err = s.GetCompatToolName("proton_513")

	assert.Empty(t, err)
	assert.Equal(t, "Proton 5.13-6", name)
}
