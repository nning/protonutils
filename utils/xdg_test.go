package utils

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_XDGDirs(t *testing.T) {
	t.Setenv("HOME", "/home/test")

	t.Setenv("XDG_CONFIG_HOME", "/xdg/config")
	t.Setenv("XDG_CACHE_HOME", "/xdg/cache")
	assert.Equal(t, "/xdg/config/protonutils", GetConfigDir())
	assert.Equal(t, "/xdg/cache/protonutils", GetCacheDir())
	assert.Equal(t, "/xdg/config/protonutils.yml", GetLegacyConfigFile())

	// Empty or relative values are ignored (XDG Base Directory Specification)
	for _, v := range []string{"", "relative/path"} {
		t.Setenv("XDG_CONFIG_HOME", v)
		t.Setenv("XDG_CACHE_HOME", v)
		assert.Equal(t, "/home/test/.config/protonutils", GetConfigDir())
		assert.Equal(t, "/home/test/.cache/protonutils", GetCacheDir())
		assert.Equal(t, "", GetLegacyConfigFile())
	}
}

func Test_ExpandHome(t *testing.T) {
	t.Setenv("HOME", "/home/test")

	assert.Equal(t, "/home/test", ExpandHome("~"))
	assert.Equal(t, filepath.Join("/home/test", "a/b"), ExpandHome("~/a/b"))
	assert.Equal(t, "~user/a", ExpandHome("~user/a"))
	assert.Equal(t, "/abs", ExpandHome("/abs"))
}

func Test_AskYesOrNo(t *testing.T) {
	for input, expected := range map[string]bool{"y\n": true, "Yes\n": true, "n\n": false, "\n": false, "": false, "y": true} {
		var out strings.Builder
		ok, err := AskYesOrNo(strings.NewReader(input), &out, "Question?")
		assert.NoError(t, err)
		assert.Equal(t, expected, ok, input)
		assert.Contains(t, out.String(), "Question? [y/N]")
	}
}
