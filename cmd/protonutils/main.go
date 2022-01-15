package main

import (
	"fmt"
	"os"

	"github.com/nning/protonutils/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	log "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:   "protonutils",
	Short: "Utilities for using the Proton compatibility tool",
	Long:  "protonutils is a CLI tool that provides different utilities to make using the Proton compatibility tool more easy.",
}

var manDir string

var cfg config.Config

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

func init() {
	rootCmd.Flags().StringVarP(&manDir, "generate-man", "m", "", "")
	rootCmd.Flags().MarkHidden("generate-man")
	rootCmd.MarkFlagDirname("generate-man")
}

func main() {
	rootCmd.ParseFlags(os.Args)

	if manDir != "" {
		header := &doc.GenManHeader{
			Title:   "PROTONUTILS",
			Section: "1",
		}
		err := os.MkdirAll(manDir, 0700)
		exitOnError(err)

		err = doc.GenManTree(rootCmd, header, manDir)
		exitOnError(err)

		return
	}

	if _, debug := os.LookupEnv("DEBUG"); debug {
		log.SetLevel(log.DebugLevel)
	}

	c, err := config.New()
	if err == nil {
		cfg = *c
		user = c.User
	}

	rootCmd.Execute()
}
