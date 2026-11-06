package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const gameSearchHelp = "This includes games that either have an explicit Proton/CompatTool mapping or have been started with Proton at least once. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively; exact matches are preferred."

type pathFunc func(idOrName string) (string, string, error)

// newPathCommands returns a parent command named name with subcommands "path"
// and "open" for the directory returned by f
func newPathCommands(name string, f pathFunc) *cobra.Command {
	parent := &cobra.Command{
		Use:   name,
		Short: "Commands regarding " + name + " directory for game",
	}

	pathCmd := &cobra.Command{
		Use:   "path [flags] <game>",
		Short: "Print " + name + " directory path for game",
		Long:  "Print " + name + " directory path for game. " + gameSearchHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePath(cmd, args, f)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), p)
			return nil
		},
	}

	openCmd := &cobra.Command{
		Use:   "open [flags] <game>",
		Short: "Open " + name + " directory for game",
		Long:  "Open " + name + " directory for game. " + gameSearchHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePath(cmd, args, f)
			if err != nil {
				return err
			}

			out, err := exec.Command("xdg-open", p).CombinedOutput()
			if err != nil {
				return fmt.Errorf("xdg-open %v: %w: %s", p, err, strings.TrimSpace(string(out)))
			}

			return nil
		},
	}

	for _, c := range []*cobra.Command{pathCmd, openCmd} {
		addUserFlags(c)
		addVerboseFlag(c)
		parent.AddCommand(c)
	}

	return parent
}

func resolvePath(cmd *cobra.Command, args []string, f pathFunc) (string, error) {
	// Populates app ID/name cache which is used for search
	err := s.ReadCompatTools()
	if err != nil {
		return "", err
	}

	p, name, err := f(strings.Join(args, " "))
	if err != nil {
		return "", err
	}

	if verbose {
		fmt.Fprintln(cmd.OutOrStdout(), name)
	}

	return p, nil
}

func init() {
	compatdataCmd := newPathCommands("compatdata", func(idOrName string) (string, string, error) {
		return s.GetCompatdataPath(idOrName)
	})
	compatdataCmd.AddCommand(compatdataCleanCmd)
	rootCmd.AddCommand(compatdataCmd)

	rootCmd.AddCommand(newPathCommands("installdir", func(idOrName string) (string, string, error) {
		return s.GetInstalldirPath(idOrName)
	}))
}
