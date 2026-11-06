package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nning/protonutils/steam"
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
	Args:  cobra.NoArgs,
	RunE:  compatToolList,
}

var compatToolSetCmd = &cobra.Command{
	Use:   "set [flags] <game> <version>",
	Short: "Set compatibility tool version for game",
	Long:  "Set compatibility tool version for game. Game search string can be app ID, game name, or prefix of game name. It is matched case-insensitively, exact matches are preferred. Version parameters have to be version IDs. See `compattool list` for list of possible options. If \"default\" is used as version, explicit mapping is removed and game uses default compatibility tool.",
	Args:  cobra.ExactArgs(2),
	RunE:  compatToolSet,
}

var compatToolMigrateCmd = &cobra.Command{
	Use:   "migrate [flags] <fromVersion> <toVersion>",
	Short: "Migrate compatibility tool version mappings from one version to another",
	Long:  "Migrate compatibility tool version mappings of installed games from one version to another. Version parameters have to be version IDs. See `compattool list` for list of possible options. If \"default\" is used as toVersion, explicit mappings are removed and games use default compatibility tool.",
	Args:  cobra.ExactArgs(2),
	RunE:  compatToolMigrate,
}

var compatToolCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete unused compatibility tools",
	Long:  "Delete compatibility tools in compatibilitytools.d that are neither used by an installed game nor configured as default.",
	Args:  cobra.NoArgs,
	RunE:  compatToolClean,
}

func init() {
	rootCmd.AddCommand(compatToolCmd)

	compatToolCmd.AddCommand(compatToolListCmd)
	addUserFlags(compatToolListCmd)

	compatToolCmd.AddCommand(compatToolSetCmd)
	addUserFlags(compatToolSetCmd)
	addYesFlag(compatToolSetCmd)

	compatToolCmd.AddCommand(compatToolMigrateCmd)
	addUserFlags(compatToolMigrateCmd)
	addYesFlag(compatToolMigrateCmd)
	compatToolMigrateCmd.Flags().BoolVarP(&remove, "remove", "r", false, "Remove fromVersion after migration (only custom tools in compatibilitytools.d)")

	compatToolCmd.AddCommand(compatToolCleanCmd)
	addUserFlags(compatToolCleanCmd)
	addYesFlag(compatToolCleanCmd)
}

func validateVersion(v string) error {
	if strings.HasPrefix(v, "proton_") || s.CompatTools.IsValid(v) || s.CustomTools.Get(v) != nil {
		return nil
	}

	return fmt.Errorf("Invalid version: %v", v)
}

func displayVersion(v string) string {
	if v == "" {
		return "default"
	}

	return v
}

// warnSteamRunning warns that Steam overwrites changes to config.vdf on exit
func warnSteamRunning(w io.Writer) {
	if isRunning, _ := s.IsRunning(); isRunning {
		fmt.Fprintf(w, "WARNING: Steam client seems to be running, changes will be overwritten on exit!\n\n")
	} else if steam.IsFlatpak() {
		fmt.Fprintf(w, "WARNING: Make sure Steam is closed, otherwise changes will be overwritten on exit!\n\n")
	}
}

func compatToolList(cmd *cobra.Command, args []string) error {
	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	for _, toolID := range s.CompatTools.Sort() {
		tool := s.CompatTools[toolID]

		if tool.IsCustom && tool.IsInstalled(s) {
			fmt.Fprintln(out, tool.ID)
			continue
		}

		if tool.Name == "" || tool.Games.CountInstalled() == 0 {
			continue
		}

		if !tool.IsCustom && !tool.IsDefault {
			fmt.Fprintln(out, tool.Name, "["+tool.ID+"]")
		} else {
			fmt.Fprintln(out, tool.Name)
		}
	}

	return nil
}

func compatToolSet(cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	newVersion := args[1]
	out := cmd.OutOrStdout()

	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	id, name, err := s.GetAppIDAndName(idOrName)
	if err != nil {
		return err
	}

	if newVersion == "default" {
		newVersion = ""
	} else if err := validateVersion(newVersion); err != nil {
		return err
	}

	oldName := "default"
	if oldVersion := s.GetGameVersion(id); oldVersion != nil {
		if oldVersion.ID == newVersion || (newVersion != "" && oldVersion.Name == newVersion) {
			fmt.Fprintf(out, "%v is already using %v\n", name, displayVersion(newVersion))
			return nil
		}

		if !oldVersion.IsDefault {
			oldName = oldVersion.Name
		}
	} else if newVersion == "" {
		fmt.Fprintf(out, "%v is already using default\n", name)
		return nil
	}

	warnSteamRunning(out)
	fmt.Fprintln(out, "App ID: ", id)
	fmt.Fprintln(out, "Name:   ", name)
	fmt.Fprintln(out)
	fmt.Fprintln(out, oldName, "->", displayVersion(newVersion))
	fmt.Fprintln(out)

	ok, err := confirm(cmd, "Really update?")
	if !ok || err != nil {
		return err
	}

	ctm := s.CompatToolMapping

	err = ctm.Update(id, newVersion)
	if err != nil {
		return err
	}

	err = ctm.Save()
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Done")
	return nil
}

func compatToolMigrate(cmd *cobra.Command, args []string) error {
	fromVersion := args[0]
	toVersion := args[1]
	out := cmd.OutOrStdout()

	if fromVersion == "default" {
		return errors.New("fromVersion cannot be default")
	}

	if toVersion == "default" {
		toVersion = ""
	}

	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	err = validateVersion(fromVersion)
	if err != nil {
		return err
	}

	if toVersion != "" {
		err = validateVersion(toVersion)
		if err != nil {
			return err
		}
	}

	if fromVersion == toVersion {
		return errors.New("fromVersion and toVersion are identical")
	}

	version := s.CompatTools[fromVersion]
	if version == nil || version.Games.CountInstalled() == 0 {
		return fmt.Errorf("No installed games for %v", fromVersion)
	}

	var removeDir string
	if remove {
		removeDir, err = removableToolDir(fromVersion)
		if err != nil {
			return err
		}
	}

	warnSteamRunning(out)
	fmt.Fprintf(out, "%v -> %v\n\n", fromVersion, displayVersion(toVersion))

	for _, id := range version.Games.Sort() {
		if game := version.Games[id]; game.IsInstalled {
			fmt.Fprintln(out, "  * "+game.Name)
		}
	}
	fmt.Fprintln(out)

	if removeDir != "" {
		fmt.Fprintf(out, "%v will be deleted afterwards.\n\n", removeDir)
	}

	ok, err := confirm(cmd, "Really update?")
	if !ok || err != nil {
		return err
	}

	ctm := s.CompatToolMapping

	for _, game := range version.Games {
		if !game.IsInstalled {
			continue
		}

		err = ctm.Update(game.ID, toVersion)
		if err != nil {
			return err
		}
	}

	err = ctm.Save()
	if err != nil {
		return err
	}

	if removeDir != "" {
		err = os.RemoveAll(removeDir)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(out, "Done")
	return nil
}

// removableToolDir returns the directory of custom tool id if it can be
// removed after migrating all installed games to another version
func removableToolDir(id string) (string, error) {
	tool := s.CustomTools.Get(id)
	if tool == nil {
		return "", fmt.Errorf("%v is not a custom compatibility tool in %v, cannot remove it", id, s.GetCompatibilityToolsDir())
	}

	if s.CompatToolMapping.GetDefaultID() == id {
		return "", fmt.Errorf("%v is configured as default compatibility tool, cannot remove it", id)
	}

	for otherID, other := range s.CustomTools {
		if otherID == id || other.Dir != tool.Dir {
			continue
		}

		if s.UsedToolIDs()[otherID] {
			return "", fmt.Errorf("%v also contains %v, which is in use, cannot remove it", tool.Dir, otherID)
		}
	}

	return tool.Dir, nil
}

func compatToolClean(cmd *cobra.Command, args []string) error {
	err := s.ReadCompatTools()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	toDelete := s.UnusedCustomToolDirs()

	// Leftovers of interrupted downloads
	leftovers, _ := filepath.Glob(filepath.Join(s.GetCompatibilityToolsDir(), tmpDirPrefix+"*"))
	toDelete = append(toDelete, leftovers...)

	if len(toDelete) == 0 {
		fmt.Fprintln(out, "No unused compatibility tool found")
		return nil
	}

	fmt.Fprintln(out, "Unused versions found:")
	for _, dir := range toDelete {
		fmt.Fprintln(out, "  * "+filepath.Base(dir))
	}
	fmt.Fprintln(out)

	ok, err := confirm(cmd, "Really delete?")
	if !ok || err != nil {
		return err
	}

	var errs []error
	for _, dir := range toDelete {
		errs = append(errs, os.RemoveAll(dir))
	}

	err = errors.Join(errs...)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Done")
	return nil
}
