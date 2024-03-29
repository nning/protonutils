package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [flags] [key] [value]",
	Short: "Get or set configuration options",
	Run:   configGetOrSet,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func configGetOrSet(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(&cfg)
		return
	}

	key := args[0]

	m := map[string]interface{}{
		"user":                     cfg.User,
		"steam_root":               cfg.SteamRoot,
		"enable_viewed_steam_play": cfg.EnableViewedSteamPlay,
	}

	if len(args) == 1 {
		fmt.Println(m[key])
	} else if len(args) >= 2 {
		switch key {
		case "user":
			cfg.User = args[1]
		case "steam_root":
			cfg.SteamRoot = args[1]
		case "enable_viewed_steam_play":
			cfg.EnableViewedSteamPlay, _ = strconv.ParseBool(args[1])
		}
		err := cfg.Save()
		exitOnError(err)
	}
}
