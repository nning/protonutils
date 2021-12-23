package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/benlubar/vdf"
	"github.com/nning/protonutils/steam"
	"github.com/nning/protonutils/utils"
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

func lookup(n *vdf.Node, x []string) (*vdf.Node, error) {
	y := n

	for _, key := range x {
		y = y.FirstByName(key)
		if y == nil {
			return nil, &steam.KeyNotFoundError{Name: key}
		}
	}

	return y, nil
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
}

func compatToolSet(cmd *cobra.Command, args []string) {
	idOrName := args[0]
	newVersion := args[1]

	s, err := steam.New(user, cfg.SteamRoot, ignoreCache)
	exitOnError(err)

	info, err := s.GetGameInfo(idOrName)
	exitOnError(err)

	oldVersion := s.GetGameVersion(info.ID)

	isValidVersion, err := s.IsValidVersion(newVersion)
	if err != nil || !isValidVersion {
		exitOnError(fmt.Errorf("Invalid version: %v", newVersion))
	}

	fmt.Println("App ID: ", info.ID)
	fmt.Println("Name:   ", info.Name)
	fmt.Println()
	fmt.Println(oldVersion, "->", newVersion)
	fmt.Println()

	if !yes {
		isOK, err := utils.AskYesOrNo("Really update?")
		exitOnError(err)

		if !isOK {
			fmt.Println("Aborted")
			return
		}
	}

	p := path.Join(s.Root, "config", "config.vdf")
	fmt.Println(p)

	in, err := ioutil.ReadFile(p)
	exitOnError(err)

	var n vdf.Node
	err = n.UnmarshalText(in)

	key := []string{"Software", "Valve", "Steam", "CompatToolMapping"}
	x, err := lookup(&n, key)
	exitOnError(err)

	y, err := lookup(x, []string{info.ID, "name"})
	_, isKeyNotFoundError := err.(*steam.KeyNotFoundError)

	if isKeyNotFoundError {
		var n0 vdf.Node
		n0.SetName(info.ID)

		var n1 vdf.Node
		n1.SetName("name")
		n1.SetString(newVersion)

		var n2 vdf.Node
		n2.SetName("config")
		n2.SetString("")

		var n3 vdf.Node
		n3.SetName("Priority")
		n3.SetString("250")

		n0.Append(&n1)
		n0.Append(&n2)
		n0.Append(&n3)

		x.Append(&n0)
	} else if err != nil {
		exitOnError(err)
	} else {
		y.SetString(newVersion)
	}

	out, err := n.MarshalText()
	ioutil.WriteFile(p, out, 0600)
}
