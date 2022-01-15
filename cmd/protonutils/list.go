package main

import (
	"encoding/json"
	"fmt"

	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List games by runtime",
	Long:  "List games by configured Proton/CompatTool version. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once.",
	Run:   list,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List both installed and non-installed games")
	listCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	listCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output JSON (implies -a and -i)")
	listCmd.Flags().BoolVarP(&showAppID, "show-id", "i", false, "Show app ID")
	listCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
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

func list(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	err = s.ReadCompatToolVersions()
	exitOnError(err)

	if !jsonOutput {
		for _, version := range s.CompatToolVersions.Sort() {
			games := s.CompatToolVersions[version].Games
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
