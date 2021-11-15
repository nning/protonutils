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

func getPath(id string) (string, error) {
	s, err := steam.New("", ignoreCache)
	if err != nil {
		return "", err
	}

	p, err := s.GetCompatdataPath(id)
	if err != nil {
		return "", err
	}

	return path.Join(p, "steamapps", "compatdata", id), nil
}

func compatdataPath(cmd *cobra.Command, args []string) {
	p, err := getPath(args[0])
	exitOnError(err)

	fmt.Println(p)
}

func compatdataOpen(cmd *cobra.Command, args []string) {
	p, err := getPath(args[0])
	exitOnError(err)

	_, err = exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
