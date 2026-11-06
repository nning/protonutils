package main

import (
	"fmt"

	"github.com/nning/protonutils/utils"
	"github.com/spf13/cobra"
)

// Version is set during build and used in output on -v
var Version string

// Buildtime is set during build and used in output on -v
var Buildtime string

var versionCmd = &cobra.Command{
	Use:         "version",
	Short:       "Show version",
	Args:        cobra.NoArgs,
	Annotations: map[string]string{steamAnnotation: steamNone},
	Run:         version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	url := "https://github.com/nning/protonutils/tree/" + Version
	f := "%-12v  %v\n"
	out := cmd.OutOrStdout()

	fmt.Fprintln(out, Version)
	fmt.Fprintf(out, "\n"+f, "Build time:", Buildtime)
	fmt.Fprintf(out, f, "Code URL:", url)
	fmt.Fprintf(out, f, "Config dir:", utils.GetConfigDir())
	fmt.Fprintf(out, f, "Cache dir:", utils.GetCacheDir())
}
