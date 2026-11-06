package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var appidCmd = &cobra.Command{
	Use:   "appid [flags] <game>",
	Short: "Search for app ID of installed game",
	Long:  "Search for app ID of installed game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively. Multiple matches are possible.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  appid,
}

func init() {
	rootCmd.AddCommand(appidCmd)
	addUserFlags(appidCmd)
}

func appid(cmd *cobra.Command, args []string) error {
	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	results := s.GetAppIDAndNames(strings.Join(args, " "))
	if len(results) == 0 {
		return errors.New("App ID could not be found")
	}

	for _, result := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "%10v  %v\n", result.ID, result.Name)
	}

	return nil
}
