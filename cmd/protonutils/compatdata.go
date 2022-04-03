package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/nning/protonutils/steam"
	"github.com/nning/protonutils/utils"
	"github.com/spf13/cobra"
)

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

var compatdataCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused compatdata directories",
	Long:  "Clean leftover compatdata directories of previously installed and now uninstalled games",
	Args:  cobra.MinimumNArgs(0),
	Run:   compatdataClean,
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

	compatdataCmd.AddCommand(compatdataCleanCmd)
	compatdataCleanCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
	compatdataCleanCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	compatdataCleanCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Do not ask")
}

func compatdataPath(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	p, n, err := s.GetCompatdataPath(strings.Join(args, " "))
	exitOnAmbiguousNameError(cmd, args, err)

	if verbose {
		fmt.Println(n)
	}

	fmt.Println(p)
}

func compatdataOpen(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	p, n, err := s.GetCompatdataPath(strings.Join(args, " "))
	exitOnAmbiguousNameError(cmd, args, err)

	if verbose {
		fmt.Println(n)
	}

	_, err = exec.Command("xdg-open", p).Output()
	exitOnError(err)
}

func compatdataClean(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	err = s.ReadCompatTools()
	exitOnError(err)

	fmt.Println("Calculating unused compatdata directory sizes...")

	type entry struct {
		name string
		path string
		size uint64
	}
	var games []entry
	var total uint64

	for _, tool := range s.CompatTools {
		for _, game := range tool.Games {
			if game.IsInstalled {
				continue
			}

			p := s.SearchCompatdataPath(game.ID)
			if p == "" {
				continue
			}

			size, err := utils.DirSize(p)
			exitOnError(err)

			games = append(games, entry{game.Name, p, size})
		}
	}

	if len(games) > 0 {
		fmt.Println()
	} else {
		fmt.Println("No unused compatdata directories found!")
		os.Exit(0)
	}

	for _, entry := range games {
		total += entry.size
		fmt.Printf("%10v  %v\n", humanize.Bytes(entry.size), entry.name)
	}

	fmt.Printf("\nTotal size: %v\n", humanize.Bytes(total))
	fmt.Println("WARNING: Backup save game data for games without Steam Cloud support!")

	if !yes {
		isOK, err := utils.AskYesOrNo("Do you want to delete compatdata directories?")
		exitOnError(err)

		if !isOK {
			fmt.Println("Aborted")
			return
		}
	}

	for _, entry := range games {
		os.RemoveAll(entry.path)
	}

	fmt.Println("Done")
}
