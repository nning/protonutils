package main

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var installdirCmd = &cobra.Command{
	Use:   "installdir",
	Short: "Commands regarding installdir directory for game",
}

var installdirPathCmd = &cobra.Command{
	Use:   "path [flags] <game>",
	Short: "Print installdir directory path for game",
	Args:  cobra.MinimumNArgs(1),
	Run:   installdirPath,
}

var installdirOpenCmd = &cobra.Command{
	Use:   "open [flags] <game>",
	Short: "Open installdir directory for game",
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

func getInstalldirPath(idOrName string) (string, string) {
	s, err := steam.New(user, ignoreCache)
	exitOnError(err)

	info, err := s.GetLibraryPath(idOrName)
	exitOnError(err)

	installdir, err := s.FindInstallDirInAppInfo(info.ID)
	exitOnError(err)

	return path.Join(info.LibraryPath, "steamapps", "common", installdir), info.Name
}

func installdirPath(cmd *cobra.Command, args []string) {
	p, n := getInstalldirPath(args[0])

	if verbose {
		fmt.Println(n)
	}

	fmt.Println(p)
}

func installdirOpen(cmd *cobra.Command, args []string) {
	p, n := getInstalldirPath(args[0])

	if verbose {
		fmt.Println(n)
	}

	_, err := exec.Command("xdg-open", p).Output()
	exitOnError(err)
}
