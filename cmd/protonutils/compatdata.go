package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

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

func getPath(idOrName string) (string, string) {
	s, err := steam.New(user, ignoreCache)
	exitOnError(err)

	p, err := s.GetLibraryPath(idOrName)
	exitOnError(err)

	var n string

	if p == "" {
		err = s.ReadCompatToolVersions()
		exitOnError(err)

		var id string
		for _, games := range s.CompatToolVersions {
			for name, game := range games {
				a := strings.ToLower(name)
				b := strings.ToLower(idOrName)

				if a == b || strings.HasPrefix(a, b) && game.IsInstalled {
					id = game.ID
					n = name
					break
				}
			}
		}

		p, err = s.GetLibraryPath(id)
		exitOnError(err)

		idOrName = id

		if id == "" || p == "" {
			fmt.Fprintln(os.Stderr, "App ID or compatdata path not found")
			os.Exit(1)
		}
	}

	return path.Join(p, "steamapps", "compatdata", idOrName), n
}

func compatdataPath(cmd *cobra.Command, args []string) {
	p, n := getPath(args[0])

	if verbose {
		fmt.Println(n)
	}

	fmt.Println(p)
}

func compatdataOpen(cmd *cobra.Command, args []string) {
	p, n := getPath(args[0])

	if verbose {
		fmt.Println(n)
	}

	_, err := exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
