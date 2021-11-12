package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "protonutils",
	Short: "Utilities for using the Proton compatibility tool",
	Long:  "protonutils is a CLI tool that provides different utilities to make using the Proton compatibility tool more easy.",
}

func exitOnError(e error, a ...string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}

	if len(a) > 0 {
		fmt.Fprintln(os.Stderr, a)
		os.Exit(1)
	}
}

func main() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
