package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nning/list_proton_versions/steam"
)

func exitOnError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func countVisibleGames(games steam.Games) int {
	i := 0

	for _, game := range games {
		if game.IsInstalled {
			i++
		}
	}

	return i
}

func main() {
	var all bool
	var ignoreCache bool
	var jsonOutput bool
	var showAppID bool
	var user string

	flag.BoolVar(&all, "a", false, "List both installed and non-installed games")
	flag.BoolVar(&ignoreCache, "c", false, "Ignore app ID/name cache")
	flag.BoolVar(&jsonOutput, "j", false, "Output JSON (implies -a and -i)")
	flag.BoolVar(&showAppID, "i", false, "Show app ID")
	flag.StringVar(&user, "u", "", "Steam user name (or SteamID3)")
	flag.Parse()

	s, err := steam.New(user, ignoreCache)
	exitOnError(err)

	err = s.ReadCompatToolVersions(user)
	exitOnError(err)

	if !jsonOutput {
		for _, version := range s.CompatToolVersions.Sort() {
			games := s.CompatToolVersions[version]
			if !all && countVisibleGames(games) == 0 {
				continue
			}

			fmt.Println(version)

			for _, game := range games.Sort() {
				if all || games[game].IsInstalled {
					fmt.Print("\t" + game)
					if showAppID {
						fmt.Print(" (" + games[game].ID + ")")
					}
					if !games[game].IsInstalled {
						fmt.Print(" [NOT INSTALLED]")
					}
					if games[game].IsShortcut {
						fmt.Print(" [SHORTCUT]")
					}
					fmt.Println()
				}
			}

			fmt.Println()
		}
	} else {
		j, err := json.MarshalIndent(s.CompatToolVersions, "", "  ")
		exitOnError(err)
		fmt.Println(string(j))
	}

	err = s.SaveCache()
	exitOnError(err)
}
