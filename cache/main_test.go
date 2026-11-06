package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Cache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	c, err := New("test", -1)
	require.NoError(t, err)

	c.Add("1", "One", true)
	c.Add("2", "", false)
	require.NoError(t, c.Write())

	c, err = New("test", -1)
	require.NoError(t, err)

	name, valid := c.Get("1")
	assert.Equal(t, "One", name)
	assert.True(t, valid)

	v, ok := c.Lookup("2")
	assert.True(t, ok)
	assert.False(t, v.Valid)

	_, ok = c.Lookup("3")
	assert.False(t, ok)

	// maxAge 0 ignores persisted entries but returns new ones
	c, err = New("test", 0)
	require.NoError(t, err)
	_, ok = c.Lookup("1")
	assert.False(t, ok)
	c.Add("4", "Four", true)
	name, _ = c.Get("4")
	assert.Equal(t, "Four", name)
}

func Test_CacheExpiry(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	old := map[string]Value{"1": {"One", true, time.Now().Add(-2 * time.Hour).UnixMicro()}}
	b, _ := json.Marshal(old)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "protonutils"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "protonutils", "test.json"), b, 0600))

	c, err := New("test", time.Hour)
	require.NoError(t, err)
	_, ok := c.Lookup("1")
	assert.False(t, ok)

	c, err = New("test", 3*time.Hour)
	require.NoError(t, err)
	_, ok = c.Lookup("1")
	assert.True(t, ok)
}

func Test_CacheCorruptFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "protonutils"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "protonutils", "test.json"), []byte("{broken"), 0600))

	c, err := New("test", -1)
	require.NoError(t, err)
	assert.Empty(t, c.Dump())
}

func Test_CacheMkdirError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	require.NoError(t, os.WriteFile(file, nil, 0600))
	t.Setenv("XDG_CACHE_HOME", file)

	_, err := New("test", -1)
	assert.Error(t, err)
}
