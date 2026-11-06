package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/nning/protonutils/config"
	"github.com/nning/protonutils/steam"
	"github.com/nning/protonutils/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Annotation key that controls how much Steam data a command needs
const steamAnnotation = "protonutils/steam"

// Values for steamAnnotation (default: full initialization)
const (
	steamNone = "none" // no Steam data needed
	steamRoot = "root" // only Steam root directory needed
)

var rootCmd = &cobra.Command{
	Use:               "protonutils",
	Short:             "Utilities for using the Proton compatibility tool",
	Long:              "protonutils is a CLI tool that provides different utilities to make using the Proton compatibility tool more easy.",
	Args:              cobra.NoArgs,
	SilenceUsage:      true,
	SilenceErrors:     true,
	RunE:              rootRun,
	PersistentPreRunE: initialize,
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if s != nil {
			return s.SaveCache()
		}
		return nil
	},
}

var manDir string
var cfg *config.Config
var s *steam.Steam
var root string

// askYesOrNo is a variable to allow replacing it in tests
var askYesOrNo = utils.AskYesOrNo

func init() {
	rootCmd.Flags().StringVarP(&manDir, "generate-man", "m", "", "")
	cobra.CheckErr(rootCmd.Flags().MarkHidden("generate-man"))
	cobra.CheckErr(rootCmd.MarkFlagDirname("generate-man"))
}

func rootRun(cmd *cobra.Command, args []string) error {
	if manDir == "" {
		return cmd.Help()
	}

	header := &doc.GenManHeader{
		Title:   "PROTONUTILS",
		Section: "1",
	}

	err := os.MkdirAll(manDir, 0700)
	if err != nil {
		return err
	}

	return doc.GenManTree(cmd.Root(), header, manDir)
}

// steamRequirement returns how much Steam data cmd needs
func steamRequirement(cmd *cobra.Command) string {
	if !cmd.HasParent() {
		return steamNone
	}

	for c := cmd; c != nil; c = c.Parent() {
		switch c.Name() {
		case "help", "completion", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			return steamNone
		}

		if v, ok := c.Annotations[steamAnnotation]; ok {
			return v
		}
	}

	return ""
}

func initialize(cmd *cobra.Command, args []string) error {
	if _, debug := os.LookupEnv("DEBUG"); debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	var err error
	cfg, err = config.New()
	if err != nil {
		return err
	}

	switch steamRequirement(cmd) {
	case steamNone:
		return nil
	case steamRoot:
		// Migrating mappings needs full Steam data
		if f := cmd.Flags().Lookup("migrate"); f == nil || !f.Changed {
			root, err = steam.FindRoot(cfg.SteamRoot)
			return err
		}
	}

	u := cfg.User
	if user != "" {
		u = user
	}

	s, err = steam.New(u, cfg, ignoreCache)
	if err != nil {
		return err
	}

	root = s.Root
	return nil
}

// printError prints err to w. For ambiguous game names, the matching games are
// listed.
func printError(w io.Writer, err error) {
	var ambiguous *steam.AmbiguousNameError
	if errors.As(err, &ambiguous) {
		for _, m := range ambiguous.Matches {
			fmt.Fprintf(w, "%10v  %v\n", m.ID, m.Name)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, err)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		printError(os.Stderr, err)
		os.Exit(1)
	}
}
