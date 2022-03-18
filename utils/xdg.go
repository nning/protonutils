package utils

import (
	"os"
	"os/user"
	"path"
	"strings"
)

func getXDGDir(name string) string {
	dir, isXDG := os.LookupEnv("XDG_" + strings.ToUpper(name) + "_HOME")
	if !isXDG {
		u, _ := user.Current()
		dir = path.Join(u.HomeDir, "."+name, "protonutils")
	}

	return dir
}

func GetConfigDir() string {
	return getXDGDir("config")
}

func GetCacheDir() string {
	return getXDGDir("cache")
}
