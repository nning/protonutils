package main

import (
	"fmt"

	"github.com/nning/protonutils/steam"
	"github.com/nning/protonutils/utils"
	"github.com/nning/protonutils/vdf2"
	"github.com/spf13/cobra"
)

var compatToolCmd = &cobra.Command{
	Use:   "compattool",
	Short: "Commands for management of compatibility tools",
}

var compatToolListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List compatibility tools",
	Long:  "List compatibility tools.",
	Run:   compatToolList,
}

var compatToolSetCmd = &cobra.Command{
	Use:   "set [flags] <game> <version>",
	Short: "Set compatibility tool version for game",
	Long:  "Set compatibility tool version for game. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively, first match is used.",
	Args:  cobra.MinimumNArgs(2),
	Run:   compatToolSet,
}

func init() {
	rootCmd.AddCommand(compatToolCmd)

	compatToolCmd.AddCommand(compatToolListCmd)
	compatToolListCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	compatToolListCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")

	compatToolCmd.AddCommand(compatToolSetCmd)
	compatToolSetCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	compatToolSetCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
	compatToolSetCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Do not ask")
}

func compatToolList(cmd *cobra.Command, args []string) {
	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	err = s.ReadCompatToolVersions()
	exitOnError(err)

	for _, versionName := range s.CompatToolVersions.Sort() {
		version := s.CompatToolVersions[versionName]
		games := version.Games

		for _, game := range games {
			if game.IsInstalled {
				id := ""
				if versionName != version.ID && !version.IsDefault {
					id = "[" + version.ID + "]"
				}
				fmt.Println(versionName, id)
				break
			}
		}
	}

	// TODO Add compat tools from .compatibilitytools.d
}

func compatToolSet(cmd *cobra.Command, args []string) {
	idOrName := args[0]
	newVersion := args[1]

	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	info, err := s.GetGameInfo(idOrName)
	exitOnError(err)

	oldVersion := s.GetGameVersion(info.ID)

	// TODO Get version ID if newVersion is name, only save mapping for ID

	isValidVersion, err := s.IsValidVersion(newVersion)
	if err != nil || !isValidVersion {
		exitOnError(fmt.Errorf("Invalid version: %v", newVersion))
	}

	if oldVersion.ID == newVersion || oldVersion.Name == newVersion {
		fmt.Println(fmt.Sprintf("%v is already using %v", info.Name, newVersion))
		return
	}

	fmt.Println("App ID: ", info.ID)
	fmt.Println("Name:   ", info.Name)
	fmt.Println()
	fmt.Println(oldVersion.Name, "->", newVersion)
	fmt.Println()

	if !yes {
		isOK, err := utils.AskYesOrNo("Really update?")
		exitOnError(err)

		if !isOK {
			fmt.Println("Aborted")
			return
		}
	}

	ctm, err := vdf2.GetCompatToolMapping(s.Root)
	exitOnError(err)

	err = ctm.Update(info.ID, newVersion)
	exitOnError(err)

	err = ctm.Save()
	exitOnError(err)

	s, err = steam.New(user, cfg.SteamRoot, true)
	exitOnError(err)

	err = s.ReadCompatToolVersions()
	exitOnError(err)

	fmt.Println("Done")
}
