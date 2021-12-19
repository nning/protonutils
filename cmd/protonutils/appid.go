package main

import (
	"fmt"
	"strings"

	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var appidCmd = &cobra.Command{
	Use:   "appid [flags] <game>",
	Short: "Search for app ID of installed game",
	Long:  "Search for app ID of installed game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be prefix of game name and is matched case-insensitively. Multiple matches are possible.",
	Args:  cobra.MinimumNArgs(1),
	Run:   appid,
}

func init() {
	rootCmd.AddCommand(appidCmd)
	appidCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	appidCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
}

func appid(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, ignoreCache)
	exitOnError(err)

	err = s.ReadCompatToolVersions()
	exitOnError(err)

	data := s.AppidCache.Dump()

	for id, value := range data {
		a := strings.ToLower(value.Name)
		b := strings.ToLower(args[0])

		if a == b || strings.HasPrefix(a, b) {
			fmt.Println(id, " ", value.Name)
		}
	}
}
