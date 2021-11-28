package main

import (
	"fmt"
	"strings"

	"github.com/nning/protonutils/cache"
	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var appidCmd = &cobra.Command{
	Use:   "appid",
	Short: "Search for app ID of installed game",
	Args:  cobra.MinimumNArgs(1),
	Run:   appid,
}

func init() {
	rootCmd.AddCommand(appidCmd)
}

func appid(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, ignoreCache)
	exitOnError(err)

	err = s.ReadCompatToolVersions()
	exitOnError(err)

	c, err := cache.New("steam-appids", -1)
	exitOnError(err)

	data := c.Dump()

	for id, value := range data {
		a := strings.ToLower(value.Name)
		b := strings.ToLower(args[0])

		if a == b || strings.HasPrefix(a, b) {
			fmt.Println(id, " ", value.Name)
		}
	}
}
