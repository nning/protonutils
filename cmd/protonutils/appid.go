package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var appidCmd = &cobra.Command{
	Use:   "appid [flags] <game>",
	Short: "Search for app ID of installed game",
	Long:  "Search for app ID of installed game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively. Multiple matches are possible.",
	Args:  cobra.MinimumNArgs(1),
	Run:   appid,
}

func init() {
	rootCmd.AddCommand(appidCmd)
	appidCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	appidCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
}

func appid(cmd *cobra.Command, args []string) {
	err := s.ReadCompatTools()
	exitOnError(err)

	results := s.GetAppIDAndNames(args[0])
	if len(results) == 0 {
		exitOnError(errors.New("App ID could not be found"))
	}

	for _, result := range results {
		fmt.Printf("%10v  %v\n", result[0], result[1])
	}
}
