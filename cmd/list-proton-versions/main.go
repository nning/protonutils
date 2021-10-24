package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/nning/list_proton_versions/steam"
)

func main() {
	var all bool
	var json_output bool

	flag.BoolVar(&all, "a", false, "List both installed and non-installed games")
	flag.BoolVar(&json_output, "j", false, "Output JSON (implies -a)")
	flag.Parse()

	s := steam.New()
	s.InitCompatToolVersions()

	if !json_output {
		for version, games := range s.CompatToolVersions {
			fmt.Println(version)

			for _, game := range games.Sort() {
				if all || games[game].IsInstalled {
					fmt.Print("\t" + game)
					if !games[game].IsInstalled {
						fmt.Print(" [NOT INSTALLED]")
					}
					fmt.Println()
				}
			}

			fmt.Println()
		}
	} else {
		j, _ := json.Marshal(s.CompatToolVersions)
		fmt.Println(string(j))
	}

	s.SaveCache()
}
