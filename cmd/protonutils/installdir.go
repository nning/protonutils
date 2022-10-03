package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var installdirCmd = &cobra.Command{
	Use:   "installdir",
	Short: "Commands regarding installdir directory for game",
}

var installdirPathCmd = &cobra.Command{
	Use:   "path [flags] <game>",
	Short: "Print installdir directory path for game",
	Long:  "Print installdir directory path for game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively.",
	Args:  cobra.MinimumNArgs(1),
	Run:   installdirPath,
}

var installdirOpenCmd = &cobra.Command{
	Use:   "open [flags] <game>",
	Short: "Open installdir directory for game",
	Long:  "Open installdir directory for game. This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively.",
	Args:  cobra.MinimumNArgs(1),
	Run:   installdirOpen,
}

func init() {
	rootCmd.AddCommand(installdirCmd)

	installdirCmd.AddCommand(installdirPathCmd)
	installdirPathCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
	installdirPathCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	installdirPathCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show app name")

	installdirCmd.AddCommand(installdirOpenCmd)
	installdirOpenCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
	installdirOpenCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	installdirOpenCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show app name")
}

func getInstalldirPath(idOrName string) (string, string, error) {
	id, name, err := s.GetAppIDAndName(idOrName)
	if err != nil {
		return "", "", err
	}

	p := s.LibraryConfig.GetLibraryPathByID(id)
	if p == "" {
		exitOnError(errors.New("Game not installed"))
	}

	dir, err := s.GetInstalldir(id)
	exitOnError(err)

	return path.Join(p, "steamapps", "common", dir), name, nil
}

func installdirPath(cmd *cobra.Command, args []string) {
	p, n, err := getInstalldirPath(strings.Join(args, " "))
	exitOnAmbiguousNameError(cmd, args, err)

	if verbose {
		fmt.Println(n)
	}

	fmt.Println(p)
}

func installdirOpen(cmd *cobra.Command, args []string) {
	p, n, err := getInstalldirPath(strings.Join(args, " "))
	exitOnAmbiguousNameError(cmd, args, err)

	if verbose {
		fmt.Println(n)
	}

	_, err = exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
