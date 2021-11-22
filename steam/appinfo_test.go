package steam

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getAppIDNeedle_uint32(t *testing.T) {
	needle, err := getAppIDNeedle("403640")
	assert.Empty(t, err)
	assert.NotEmpty(t, needle)

	x := bytes.Compare(needle, []byte{'a', 'p', 'p', 'i', 'd', 0, 0xb8, 0x28, 6, 0})
	assert.Equal(t, 0, x, "Needle should be correct")
}

func Test_getAppIDNeedle_uint64(t *testing.T) {
	needle, err := getAppIDNeedle("3372497361")
	assert.Empty(t, err)
	assert.NotEmpty(t, needle)

	x := bytes.Compare(needle, []byte{'a', 'p', 'p', 'i', 'd', 0, 0xd1, 0x39, 4, 0xc9})
	assert.Equal(t, 0, x, "Needle should be correct")
}

func Test_getAppInfoBuffer(t *testing.T) {
	s, err := New("", true)
	assert.Empty(t, err)

	info, buf, err := s.getAppInfoBuffer()
	assert.Empty(t, err)
	assert.NotEmpty(t, buf)
	assert.Equal(t, appInfoMagic, info.Magic, "AppInfo Magic invalid")
	assert.Equal(t, uint32(1), info.Universe, "Universe invalid")
}

func Test_FindNameInAppInfo_found(t *testing.T) {
	s, err := New("", true)
	assert.Empty(t, err)

	games := map[string]string{
		"292030":  "The Witcher 3: Wild Hunt",
		"377160":  "Fallout 4",
		"403640":  "Dishonored 2",
		"614570":  "Dishonored®: Death of the Outsider™ ",
		"826630":  "Iron Harvest",
		"1091500": "Cyberpunk 2077",
		"1151640": "Horizon Zero Dawn",
		"1174180": "Red Dead Redemption 2",
		"1328670": "Mass Effect™ Legendary Edition",
	}

	for id, n := range games {
		name, err := s.findNameInAppInfo(id)
		assert.Empty(t, err)
		assert.Equal(t, n, name)
	}
}

func Test_FindNameInAppInfo_notFound(t *testing.T) {
	s, err := New("", true)
	assert.Empty(t, err)

	name, err := s.findNameInAppInfo("386360")
	assert.True(t, errors.Is(err, io.EOF))
	assert.Equal(t, "", name)
}

func Test_FindNameInShortcuts_found(t *testing.T) {
	s, err := New("", true)
	assert.Empty(t, err)

	name, err := s.findNameInShortcuts("3228583970")
	assert.Empty(t, err)
	assert.Equal(t, "Kena - Bridge of Spirits", name)
}
