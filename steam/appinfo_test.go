package steam

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getNeedle(t *testing.T) {
	needle, err := getNeedle("403640")
	assert.Empty(t, err)
	assert.NotEmpty(t, needle)

	x := bytes.Compare(needle, []byte{'a', 'p', 'p', 'i', 'd', 0, 0xb8, 0x28, 6, 0})
	assert.Equal(t, 0, x, "Needle should be correct")
}

func Test_getAppInfoBuffer(t *testing.T) {
	s, err := New(true)
	assert.Empty(t, err)

	info, buf, err := s.getAppInfoBuffer()
	assert.Empty(t, err)
	assert.NotEmpty(t, buf)
	assert.Equal(t, AppInfoMagic, info.Magic, "AppInfo Magic invalid")
	assert.Equal(t, uint32(1), info.Universe, "Universe invalid")
}
