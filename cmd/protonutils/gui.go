package main

import (
	"github.com/nning/protonutils/gui"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Start GUI",
	Run:   guiRun,
}

func init() {
	rootCmd.AddCommand(guiCmd)
	guiCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
}

func guiRun(cmd *cobra.Command, args []string) {
	gui.Run(user, &cfg, ignoreCache)
}
