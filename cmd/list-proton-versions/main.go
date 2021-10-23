package main

import (
	"fmt"

	. "github.com/nning/list_proton_versions"
	"github.com/nning/list_proton_versions/set"
	"github.com/nning/list_proton_versions/steam"
)

type Versions map[string]set.Set

func (versions Versions) Includes(name string) bool {
	x := false

	for _, games := range versions {
		if games.Includes(name) {
			return true
		}
	}

	return x
}

func main() {
	s := steam.New()

	x, err := s.GetCompatToolMapping()
	PanicOnError(err)

	versions := make(Versions)

	for id, cfg := range x {
		if id == "0" {
			continue
		}

		v := cfg.(steam.MapLevel)["name"].(string)
		if v == "" {
			v = "Default"
		}

		name := s.GetNameWithInstallStatus(id)
		if name != "" {
			versions[v] = set.Init(versions[v])
			versions[v].Add(name)
		}
	}

	x, err = s.GetLocalConfig()
	PanicOnError(err)

	for id, cfg := range x {
		v := cfg.(steam.MapLevel)["ViewedSteamPlay"]
		if v == nil {
			continue
		}

		name := s.GetNameWithInstallStatus(id)
		if name != "" && !versions.Includes(name) {
			versions["Default"].Add(name)
		}
	}

	for version, games := range versions {
		fmt.Println(version)

		for _, game := range games.Sort() {
			fmt.Println("\t" + game)
		}

		fmt.Println()
	}

	s.SaveCache()
}
