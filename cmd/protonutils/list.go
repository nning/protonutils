package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/nning/protonutils/steam"
	"github.com/nning/protonutils/utils"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List games by runtime",
	Long:  "List games by configured Proton/CompatTool version. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once.",
	Args:  cobra.NoArgs,
	RunE:  list,
}

func init() {
	rootCmd.AddCommand(listCmd)
	addUserFlags(listCmd)
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List both installed and non-installed games")
	listCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output JSON (implies -a and -i)")
	listCmd.Flags().BoolVarP(&showAppID, "show-id", "i", false, "Show app ID")
	listCmd.Flags().BoolVarP(&showSize, "size", "s", false, "Show disk usage of compatdata and game directories (JSON: also paths)")
}

func list(cmd *cobra.Command, args []string) error {
	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	var totalCompatdata, totalInstall uint64
	if showSize {
		totalCompatdata, totalInstall = addSizes(jsonOutput || all)
	}

	if jsonOutput {
		j, err := json.MarshalIndent(s.CompatTools, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(out, string(j))
		return nil
	}

	for _, toolID := range s.CompatTools.Sort() {
		tool := s.CompatTools[toolID]
		games := tool.Games
		if !all && games.CountInstalled() == 0 {
			continue
		}

		if tool.Name == "" {
			fmt.Fprintln(out, "Default")
		} else {
			fmt.Fprintln(out, tool.Name)
		}

		for _, id := range games.Sort() {
			game := games[id]
			if !all && !game.IsInstalled {
				continue
			}

			fmt.Fprint(out, "\t"+game.Name)
			if showAppID {
				fmt.Fprint(out, " ("+game.ID+")")
			}
			if !game.IsInstalled {
				fmt.Fprint(out, " [NOT INSTALLED]")
			}
			if game.IsShortcut {
				fmt.Fprint(out, " [SHORTCUT]")
			}
			if showSize {
				fmt.Fprint(out, formatSizes(game))
			}
			fmt.Fprintln(out)
		}

		fmt.Fprintln(out)
	}

	if showSize {
		fmt.Fprintf(out, "Total: compatdata %v, games %v\n", humanize.Bytes(totalCompatdata), humanize.Bytes(totalInstall))
	}

	return nil
}

// addSizes sets paths and sizes of all (or only installed) games and returns
// total sizes
func addSizes(all bool) (compatdata, install uint64) {
	for _, tool := range s.CompatTools {
		for _, game := range tool.Games {
			if !all && !game.IsInstalled {
				continue
			}

			game.CompatdataPath, game.InstallPath = s.GetGamePaths(game.ID)

			if game.CompatdataPath != "" {
				game.CompatdataSize, _ = utils.DirSize(game.CompatdataPath)
				compatdata += game.CompatdataSize
			}

			if game.InstallPath != "" {
				game.InstallSize, _ = utils.DirSize(game.InstallPath)
				install += game.InstallSize
			}
		}
	}

	return compatdata, install
}

// formatSizes returns sizes of game directories for list output
func formatSizes(game *steam.Game) string {
	var parts []string

	if game.CompatdataPath != "" {
		parts = append(parts, "compatdata "+humanize.Bytes(game.CompatdataSize))
	}

	if game.InstallPath != "" {
		parts = append(parts, "game "+humanize.Bytes(game.InstallSize))
	}

	if len(parts) == 0 {
		return ""
	}

	return " [" + strings.Join(parts, ", ") + "]"
}
