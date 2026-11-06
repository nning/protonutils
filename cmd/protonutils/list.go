package main

import (
	"encoding/json"
	"fmt"

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
}

func list(cmd *cobra.Command, args []string) error {
	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

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
			fmt.Fprintln(out)
		}

		fmt.Fprintln(out)
	}

	return nil
}
