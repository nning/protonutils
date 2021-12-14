package main

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var compatdataCmd = &cobra.Command{
	Use:   "compatdata",
	Short: "Commands regarding compatdata directory for game",
}

var compatdataPathCmd = &cobra.Command{
	Use:   "path [flags] <game>",
	Short: "Print compatdata directory path for game",
	Args:  cobra.MinimumNArgs(1),
	Run:   compatdataPath,
}

var compatdataOpenCmd = &cobra.Command{
	Use:   "open [flags] <game>",
	Short: "Open compatdata directory for game",
	Args:  cobra.MinimumNArgs(1),
	Run:   compatdataOpen,
}

var verbose bool

func init() {
	rootCmd.AddCommand(compatdataCmd)

	compatdataCmd.AddCommand(compatdataPathCmd)
	compatdataPathCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
	compatdataPathCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	compatdataPathCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show app name")

	compatdataCmd.AddCommand(compatdataOpenCmd)
	compatdataOpenCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
	compatdataOpenCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	compatdataOpenCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show app name")
}

func getCompatdataPath(idOrName string) (string, string) {
	s, err := steam.New(user, ignoreCache)
	exitOnError(err)

	info, err := s.GetGameInfo(idOrName)
	exitOnError(err)

	return path.Join(info.LibraryPath, "steamapps", "compatdata", info.ID), info.Name
}

func compatdataPath(cmd *cobra.Command, args []string) {
	p, n := getCompatdataPath(args[0])

	if verbose {
		fmt.Println(n)
	}

	fmt.Println(p)
}

func compatdataOpen(cmd *cobra.Command, args []string) {
	p, n := getCompatdataPath(args[0])

	if verbose {
		fmt.Println(n)
	}

	_, err := exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
