package main

import (
	"encoding/json"
	"flag"
	"fmt"

	. "github.com/nning/list_proton_versions"
	"github.com/nning/list_proton_versions/steam"
)

func main() {
	var all bool
	var ignore_cache bool
	var json_output bool
	var user string

	flag.BoolVar(&all, "a", false, "List both installed and non-installed games")
	flag.BoolVar(&json_output, "j", false, "Output JSON (implies -a)")
	flag.BoolVar(&ignore_cache, "i", false, "Ignore app id/name cache")
	flag.StringVar(&user, "u", "", "Steam user name (or SteamID3)")
	flag.Parse()

	s, err := steam.New(!ignore_cache)
	ExitOnError(err)

	err = s.InitCompatToolVersions(user)
	ExitOnError(err)

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
		j, _ := json.MarshalIndent(s.CompatToolVersions, "", "  ")
		fmt.Println(string(j))
	}

	err = s.SaveCache()
	ExitOnError(err)
}
