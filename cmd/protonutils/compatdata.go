package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path"

	"github.com/nning/protonutils/steam2"
	"github.com/spf13/cobra"
)

type AmbiguousNameError struct{}

func (err *AmbiguousNameError) Error() string {
	return "Ambiguous name, try using app ID"
}

var compatdataCmd = &cobra.Command{
	Use:   "compatdata",
	Short: "Commands regarding compatdata directory for game",
}

var compatdataPathCmd = &cobra.Command{
	Use:   "path [flags] <game>",
	Short: "Print compatdata directory path for game",
	Long:  "Print compatdata directory path for game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively.",
	Args:  cobra.MinimumNArgs(1),
	Run:   compatdataPath,
}

var compatdataOpenCmd = &cobra.Command{
	Use:   "open [flags] <game>",
	Short: "Open compatdata directory for game",
	Long:  "Open compatdata directory for game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively.",
	Args:  cobra.MinimumNArgs(1),
	Run:   compatdataOpen,
}

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

func getCompatdataPath(idOrName string) (string, string, error) {
	s, err := steam2.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	idAndNames := s.GetAppIDAndNames(idOrName)

	l := len(idAndNames)
	if l == 0 {
		return "", "", errors.New("App ID or name not found")
	} else if l > 1 {
		return "", "", &AmbiguousNameError{}
	}

	id := idAndNames[0][0]
	name := idAndNames[0][1]

	p := s.LibraryConfig.GetLibraryPathByID(id)
	if p == "" {
		exitOnError(errors.New("Game not installed"))
	}

	return path.Join(p, "steamapps", "compatdata", id), name, nil
}

func checkError(cmd *cobra.Command, args []string, err error) {
	if err != nil {
		if _, isAmbiguous := err.(*AmbiguousNameError); isAmbiguous {
			appid(cmd, args)
			fmt.Println()
		}

		exitOnError(err)
	}
}

func compatdataPath(cmd *cobra.Command, args []string) {
	p, n, err := getCompatdataPath(args[0])
	checkError(cmd, args, err)

	if verbose {
		fmt.Println(n)
	}

	fmt.Println(p)
}

func compatdataOpen(cmd *cobra.Command, args []string) {
	p, n, err := getCompatdataPath(args[0])
	checkError(cmd, args, err)

	if verbose {
		fmt.Println(n)
	}

	_, err = exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
