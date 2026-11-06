package main

import (
	"fmt"

	"github.com/nning/protonutils/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:         "config [flags] [key] [value]",
	Short:       "Get or set configuration options",
	Long:        "Get or set configuration options. Without arguments, all options are shown. Valid keys: user, steam_root (empty for auto-detection), enable_viewed_steam_play.",
	Args:        cobra.MaximumNArgs(2),
	ValidArgs:   config.Keys(),
	Annotations: map[string]string{steamAnnotation: steamNone},
	RunE:        configGetOrSet,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func configGetOrSet(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()

	if len(args) == 0 {
		fmt.Fprint(out, cfg)
		return nil
	}

	if len(args) == 1 {
		v, err := cfg.Get(args[0])
		if err != nil {
			return err
		}

		fmt.Fprintln(out, v)
		return nil
	}

	err := cfg.Set(args[0], args[1])
	if err != nil {
		return err
	}

	return cfg.Save()
}
