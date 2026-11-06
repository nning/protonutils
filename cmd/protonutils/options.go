package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var all bool
var extractOnly bool
var force bool
var ignoreCache bool
var jsonOutput bool
var keep bool
var migrate bool
var remove bool
var showAppID bool
var showSize bool
var user string
var verbose bool
var yes bool

func addUserFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name, steamID64 or account ID (SteamID3)")
	cmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
}

func addYesFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Do not ask")
}

func addVerboseFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show app name")
}

// confirm asks the user for confirmation (unless --yes is set) and prints
// "Aborted" if the answer is no
func confirm(cmd *cobra.Command, question string) (bool, error) {
	if yes {
		return true, nil
	}

	ok, err := askYesOrNo(cmd.InOrStdin(), cmd.OutOrStdout(), question)
	if err != nil {
		return false, err
	}

	if !ok {
		fmt.Fprintln(cmd.OutOrStdout(), "Aborted")
	}

	return ok, nil
}
