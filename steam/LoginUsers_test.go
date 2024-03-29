package steam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetID64(t *testing.T) {
	s, err := New("", testConfig, false)

	assert.Empty(t, err)
	assert.Equal(t, "76561198050517827", s.LoginUsers.GetID64("phects"))
	assert.Equal(t, "76561198050517828", s.LoginUsers.GetID64("phectt"))
	assert.Equal(t, "", s.LoginUsers.GetID64(""))
	assert.Equal(t, "", s.LoginUsers.GetID64("foo"))

	s, err = New("phects", testConfig, false)

	assert.Empty(t, err)
	assert.Equal(t, "76561198050517827", s.LoginUsers.GetID64("phects"))
}
