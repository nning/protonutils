package utils

import (
	"os"
	"path/filepath"
)

const appName = "protonutils"

// HomeDir returns the home directory of the current user (or "/" if it cannot
// be determined)
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "/"
	}

	return home
}

// ExpandHome replaces a leading "~/" (or a sole "~") in p with the home
// directory of the current user
func ExpandHome(p string) string {
	if p == "~" {
		return HomeDir()
	}

	if len(p) > 1 && p[0] == '~' && p[1] == '/' {
		return filepath.Join(HomeDir(), p[2:])
	}

	return p
}

// xdgBaseDir returns the XDG base directory for kind ("config" or "cache").
// According to the XDG Base Directory Specification, the environment variable
// is ignored if it is empty or not an absolute path.
func xdgBaseDir(kind, envName string) string {
	dir := os.Getenv(envName)
	if dir != "" && filepath.IsAbs(dir) {
		return dir
	}

	return filepath.Join(HomeDir(), "."+kind)
}

// GetConfigDir returns XDG config dir of protonutils
func GetConfigDir() string {
	return filepath.Join(xdgBaseDir("config", "XDG_CONFIG_HOME"), appName)
}

// GetCacheDir returns XDG cache dir of protonutils
func GetCacheDir() string {
	return filepath.Join(xdgBaseDir("cache", "XDG_CACHE_HOME"), appName)
}

// GetLegacyConfigFile returns the config file location used by protonutils
// versions that did not append "protonutils" to $XDG_CONFIG_HOME (empty string
// if XDG_CONFIG_HOME is not set)
func GetLegacyConfigFile() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" || !filepath.IsAbs(dir) {
		return ""
	}

	return filepath.Join(dir, appName+".yml")
}
