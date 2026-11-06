package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/nning/protonutils/utils"
	"github.com/spf13/cobra"
)

var compatdataCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused compatdata directories",
	Long:  "Clean leftover compatdata directories of previously installed and now uninstalled games",
	Args:  cobra.NoArgs,
	RunE:  compatdataClean,
}

func init() {
	addUserFlags(compatdataCleanCmd)
	addYesFlag(compatdataCleanCmd)
}

func compatdataClean(cmd *cobra.Command, args []string) error {
	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "Calculating unused compatdata directory sizes...")

	type entry struct {
		name string
		path string
		size uint64
	}
	var entries []entry
	var total uint64

	seen := make(map[string]bool)
	for _, tool := range s.CompatTools {
		for _, game := range tool.Games {
			if game.IsInstalled || seen[game.ID] {
				continue
			}
			seen[game.ID] = true

			p := s.SearchCompatdataPath(game.ID)
			if p == "" {
				continue
			}

			size, err := utils.DirSize(p)
			if err != nil {
				return err
			}

			entries = append(entries, entry{game.Name, p, size})
		}
	}

	if len(entries) == 0 {
		fmt.Fprintln(out, "No unused compatdata directories found!")
		return nil
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })

	fmt.Fprintln(out)
	for _, e := range entries {
		total += e.size
		fmt.Fprintf(out, "%10v  %v\n", humanize.Bytes(e.size), e.name)
	}

	fmt.Fprintf(out, "\nTotal size: %v\n", humanize.Bytes(total))
	fmt.Fprintln(out, "WARNING: Backup save game data for games without Steam Cloud support!")

	ok, err := confirm(cmd, "Do you want to delete compatdata directories?")
	if !ok || err != nil {
		return err
	}

	var errs []error
	for _, e := range entries {
		errs = append(errs, os.RemoveAll(e.path))
	}

	err = errors.Join(errs...)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Done")
	return nil
}
