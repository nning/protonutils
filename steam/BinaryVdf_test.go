package steam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetNextEntryStartByID(t *testing.T) {
	s, err := New("", testConfig, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	id := "1091500"
	pos := 10233146

	i, err := s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, id)
	assert.Empty(t, err)
	assert.Equal(t, pos, i)

	i, err = s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, "0")
	assert.Empty(t, err)
	assert.Equal(t, -1, i)

	i, err = s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, "")
	assert.Empty(t, err)
	assert.Equal(t, -1, i)

	i, err = s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, "41414141")
	assert.Empty(t, err)
	assert.Equal(t, -1, i)
}

func Test_ParseAppInfoEntry(t *testing.T) {
	s, err := New("", testConfig, false)
	assert.Empty(t, err)
	assert.NotEmpty(t, s.AppInfo)

	id := "1091500"
	name := "Cyberpunk 2077"

	i, err := s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, id)
	assert.Empty(t, err)

	n, err := ParseBinaryVdf(s.AppInfo.Bytes[i:])
	assert.Empty(t, err)
	assert.NotEmpty(t, n)
	assert.Equal(t, id, n.FirstByName("appid").String())
	assert.Equal(t, name, n.FirstByName("common").FirstByName("name").String())
	assert.Equal(t, "appinfo", n.Name())
}

func Test_GetDeckCompatibility(t *testing.T) {
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
		i, err := s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, id)
		assert.Empty(t, err)
		assert.True(t, i > 0, id)

		n, err := ParseBinaryVdf(s.AppInfo.Bytes[i:])
		assert.Empty(t, err)
		assert.NotEmpty(t, n)
		assert.Equal(t, id, n.FirstByName("appid").String())

		cn := n.FirstByName("common").FirstByName("steam_deck_compatibility")
		c := GetDeckCompatibility(cn)

		assert.Equal(t, expectedCategory, c.Category)
	}
}
