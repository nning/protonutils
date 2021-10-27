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
	var ignoreCache bool
	var jsonOutput bool
	var showAppId bool
	var user string

	flag.BoolVar(&all, "a", false, "List both installed and non-installed games")
	flag.BoolVar(&ignoreCache, "c", false, "Ignore app ID/name cache")
	flag.BoolVar(&jsonOutput, "j", false, "Output JSON (implies -a)")
	flag.BoolVar(&showAppId, "i", false, "Show app ID")
	flag.StringVar(&user, "u", "", "Steam user name (or SteamID3)")
	flag.Parse()

	s, err := steam.New(!ignoreCache)
	ExitOnError(err)

	err = s.InitCompatToolVersions(user)
	ExitOnError(err)

	if !jsonOutput {
		for version, games := range s.CompatToolVersions {
			fmt.Println(version)

			for _, game := range games.Sort() {
				if all || games[game].IsInstalled {
					fmt.Print("\t" + game)
					if showAppId {
						fmt.Print(" (" + games[game].Id + ")")
					}
					if !games[game].IsInstalled {
						fmt.Print(" [NOT INSTALLED]")
					}
					fmt.Println()
				}
			}

			fmt.Println()
		}
	} else {
		j, err := json.MarshalIndent(s.CompatToolVersions, "", "  ")
		ExitOnError(err)
		fmt.Println(string(j))
	}

	err = s.SaveCache()
	ExitOnError(err)
}
