package main

import (
	"fmt"

	. "github.com/nning/list_proton_versions"
	"github.com/nning/list_proton_versions/steam"
)

func main() {
	steam := steam.New()

	x, err := steam.GetCompatToolMapping()
	PanicOnError(err)

	versions := make(map[string][]string)

	for id, cfg := range x {
		if id == "0" {
			continue
		}

		v := cfg.(map[string]interface{})["name"].(string)

		if v == "" {
			v = "Default"
		}

		if versions[v] == nil {
			versions[v] = make([]string, 0)
		}

		name := steam.GetName(id)
		installed := steam.IsInstalled(id)
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

	steam.SaveCache()
}
