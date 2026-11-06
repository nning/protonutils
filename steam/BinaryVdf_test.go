package steam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AppInfoGet(t *testing.T) {
	requireTestRoot(t)
	s, err := New("", testConfig, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	n, err := s.AppInfo.Get("1091500")
	assert.Empty(t, err)
	assert.NotNil(t, n)

	for _, id := range []string{"0", "", "41414141"} {
		n, err = s.AppInfo.Get(id)
		assert.Empty(t, err)
		assert.Nil(t, n)
	}
}

func Test_ParseAppInfoEntry(t *testing.T) {
	requireTestRoot(t)
	s, err := New("", testConfig, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	id := "1091500"
	name := "Cyberpunk 2077"

	n, err := s.AppInfo.Get(id)
	assert.Empty(t, err)
	assert.NotEmpty(t, n)
	assert.Equal(t, id, n.FirstByName("appid").String())
	assert.Equal(t, name, n.FirstByName("common").FirstByName("name").String())
	assert.Equal(t, "appinfo", n.Name())
}

func Test_GetDeckCompatibility(t *testing.T) {
	requireTestRoot(t)
	s, err := New("", testConfig, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	tests := map[string]DeckCompatibilityCategory{
		"1091500": DeckCompatibilityUnknown,  // Cyberpunk 2077
		"1426210": DeckCompatibilityPlayable, // It Takes Two
		"292030":  DeckCompatibilityPlayable, // The Witcher 3
		"205100":  DeckCompatibilityVerified, // Dishonored
		"620":     DeckCompatibilityVerified, // Portal 2
		"1190460": DeckCompatibilityVerified, // Death Stranding
		// TODO Find example for DeckCompatibilityUnsupported
	}

	for id, expectedCategory := range tests {
		n, err := s.AppInfo.Get(id)
		assert.Empty(t, err)
		assert.NotNil(t, n, id)
		assert.Equal(t, id, n.FirstByName("appid").String())

		cn := n.FirstByName("common").FirstByName("steam_deck_compatibility")
		c := GetDeckCompatibility(cn)

		assert.Equal(t, expectedCategory, c.Category)
	}
}
