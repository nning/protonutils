package main

import (
	"fmt"
	"os"

	"github.com/nning/list_proton_versions/steam"
)

func printHelp() {
	fmt.Println(
		`list-proton-versions [options]

    -a      List both installed and non-installed games
    -h      Show this help text
`)
}

func main() {
	all := false
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-a":
			all = true
		case "-h":
			printHelp()
			os.Exit(1)
		}
	}

	s := steam.New()
	s.InitCompatToolVersions()

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

	s.SaveCache()
}
