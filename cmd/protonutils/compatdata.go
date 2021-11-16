package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/nning/protonutils/cache"
	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var compatdataCmd = &cobra.Command{
	Use:   "compatdata",
	Short: "Commands regarding compatdata directory for game",
}

var compatdataPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print compatdata directory path for game",
	Args:  cobra.MinimumNArgs(1),
	Run:   compatdataPath,
}

var compatdataOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open compatdata directory for game",
	Args:  cobra.MinimumNArgs(1),
	Run:   compatdataOpen,
}

func init() {
	rootCmd.AddCommand(compatdataCmd)
	compatdataCmd.AddCommand(compatdataPathCmd)
	compatdataCmd.AddCommand(compatdataOpenCmd)
}

func getPath(idOrName string) string {
	s, err := steam.New("", ignoreCache)
	exitOnError(err)

	p, err := s.GetCompatdataPath(idOrName)
	exitOnError(err)

	if p == "" {
		err = s.ReadCompatToolVersions()
		exitOnError(err)

		c, err := cache.New("steam-appids", false)
		exitOnError(err)

		var id2 string
		for id, value := range c.Dump() {
			a := strings.ToLower(value.Name)
			b := strings.ToLower(idOrName)

			if a == b || strings.HasPrefix(a, b) {
				id2 = id
				break
			}
		}

		if id2 == "" {
			fmt.Fprintln(os.Stderr, "App ID or compatdata path not found")
			os.Exit(1)
		}

		idOrName = id2
		p, err = s.GetCompatdataPath(idOrName)
		exitOnError(err)
	}

	return path.Join(p, "steamapps", "compatdata", idOrName)
}

func compatdataPath(cmd *cobra.Command, args []string) {
	p := getPath(args[0])
	fmt.Println(p)
}

func compatdataOpen(cmd *cobra.Command, args []string) {
	p := getPath(args[0])
	_, err := exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
