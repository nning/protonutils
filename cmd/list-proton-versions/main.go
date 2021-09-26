package main

import (
	"fmt"

	. "github.com/nning/list_proton_versions"
	"github.com/nning/list_proton_versions/steam"
)

type Versions map[string][]string

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

		// if versions[v] == nil {
		// 	versions[v] = make([]string, 0)
		// }

		name := s.GetName(id)
		installed := s.IsInstalled(id)
		ni := ""
		if !installed {
			ni = " [NOT INSTALLED]"
		}

		if name != "💩" {
			versions[v] = append(versions[v], name+ni)
		}
	}

	for version, games := range versions {
		fmt.Println(version)

		for _, game := range games {
			fmt.Println("\t" + game)
		}

		fmt.Println()
	}

	s.SaveCache()
}
