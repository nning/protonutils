package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nning/list_proton_versions/steam"
)

func printHelp() {
	fmt.Println(
		`list-proton-versions [options]

    -a      List both installed and non-installed games
    -h      Show this help text
    -j      Output JSON
`)
}

func main() {
	all := false
	json_output := false
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-a":
			all = true
		case "-h":
			printHelp()
			os.Exit(1)
		case "-j":
			json_output = true
		}
	}

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
