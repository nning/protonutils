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
	var user string

	flag.BoolVar(&all, "a", false, "List both installed and non-installed games")
	flag.BoolVar(&json_output, "j", false, "Output JSON (implies -a)")
	flag.StringVar(&user, "u", "", "Steam user ID")
	flag.Parse()

	s := steam.New()
	s.InitCompatToolVersions(user)

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
