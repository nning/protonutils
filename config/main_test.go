package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ConfigSaveLoad(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, "", cfg.SteamRoot)

	require.NoError(t, cfg.Set("user", "alice"))
	require.NoError(t, cfg.Set("enable_viewed_steam_play", "true"))
	assert.Error(t, cfg.Set("enable_viewed_steam_play", "maybe"))
	assert.Error(t, cfg.Set("foo", "bar"))
	require.NoError(t, cfg.Save())

	cfg, err = New()
	require.NoError(t, err)
	assert.Equal(t, "alice", cfg.User)
	assert.True(t, cfg.EnableViewedSteamPlay)

	v, err := cfg.Get("enable_viewed_steam_play")
	assert.NoError(t, err)
	assert.Equal(t, "true", v)

	_, err = cfg.Get("foo")
	assert.Error(t, err)
}

func Test_ConfigInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "protonutils"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "protonutils", "protonutils.yml"), []byte("user: [\n"), 0600))

	_, err := New()
	assert.ErrorContains(t, err, "protonutils.yml")
}

func Test_ConfigLegacyMigration(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	legacy := filepath.Join(dir, "protonutils.yml")
	require.NoError(t, os.WriteFile(legacy, []byte("user: bob\n"), 0600))

	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, "bob", cfg.User)
	assert.Equal(t, filepath.Join(dir, "protonutils", "protonutils.yml"), cfg.Path())
	assert.FileExists(t, cfg.Path())
	assert.NoFileExists(t, legacy)
}
