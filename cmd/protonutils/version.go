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
	f := "%-12v  %v\n"

	fmt.Println(Version)
	fmt.Printf("\n"+f, "Build time:", Buildtime)
	fmt.Printf(f, "Code URL:", url)
	fmt.Printf(f, "Config dir:", utils.GetConfigDir())
	fmt.Printf(f, "Cache dir:", utils.GetCacheDir())
}
