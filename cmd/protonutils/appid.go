package main

import (
	"fmt"

	"github.com/nning/protonutils/cache"
	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var appidCmd = &cobra.Command{
	Use:   "appid",
	Short: "Search for app ID of installed game",
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

	c, err := cache.New("steam-appids", false)
	exitOnError(err)

	data := c.Dump()

	for id, value := range data {
		if value.Name == args[0] {
			fmt.Println(id)
		}
	}
}
