package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set during build and used in output on -v
var Version string

// Buildtime is set during build and used in output on -v
var Buildtime string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run:   version,
}

var showVersion bool

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	url := "https://github.com/nning/protonutils/tree/" + Version
	fmt.Println(Version, Buildtime, url)
}
