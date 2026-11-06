package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

const geRepo = "GloriousEggroll/proton-ge-custom"
const luxtorpedaRepo = "luxtorpeda-dev/luxtorpeda"

var (
	geTagOld       = regexp.MustCompile(`^[0-9]+\.[0-9]+(rc[0-9]+)?-GE-[0-9]+(-[A-Za-z0-9]+)*$`)
	geTagNew       = regexp.MustCompile(`^GE-Proton[0-9]+-[0-9]+(-[A-Za-z0-9]+)*$`)
	luxtorpedaTag  = regexp.MustCompile(`^v[0-9]+(\.[0-9]+)*([-.][A-Za-z0-9]+)*$`)
	toolAnnotation = map[string]string{steamAnnotation: steamRoot}
)

func newToolCommands(use, name string, download func(cmd *cobra.Command, tag string) error, repo string, withExtractOnly bool) *cobra.Command {
	parent := &cobra.Command{
		Use:         use,
		Short:       "Commands for " + name,
		Annotations: toolAnnotation,
	}

	downloadCmd := &cobra.Command{
		Use:   "download [flags] <version>",
		Short: "Download and extract version specified in argument",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return download(cmd, args[0])
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Download and extract the latest " + name + " release",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tag, err := latestReleaseTag(repo)
			if err != nil {
				return err
			}

			return download(cmd, tag)
		},
	}

	downloadCmd.Flags().BoolVarP(&force, "force", "f", false, "Force download even if version exists")
	updateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force last version update")

	for _, c := range []*cobra.Command{downloadCmd, updateCmd} {
		c.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")
		if withExtractOnly {
			c.Flags().BoolVarP(&extractOnly, "extract-only", "e", false, "Do not download but extract only from existing archive")
		}

		parent.AddCommand(c)
	}

	return parent
}

func init() {
	rootCmd.AddCommand(newToolCommands("ge", "Proton-GE", geDownload, geRepo, false))
	rootCmd.AddCommand(newToolCommands("luxtorpeda", "Luxtorpeda", luxtorpedaDownload, luxtorpedaRepo, true))
}

func geDownload(cmd *cobra.Command, tag string) error {
	validOld := geTagOld.MatchString(tag)
	if !validOld && !geTagNew.MatchString(tag) {
		return fmt.Errorf("No valid GE version tag: %q", tag)
	}

	name := tag
	if validOld {
		name = "Proton-" + tag
	}

	return installTool(cmd.OutOrStdout(), toolRelease{
		candidates: []toolCandidate{{name: name, archive: name + ".tar.gz"}},
		baseURL:    githubURL + "/" + geRepo + "/releases/download/" + tag,
		checksum: func(c toolCandidate) string {
			return c.name + ".sha512sum"
		},
	})
}

func luxtorpedaDownload(cmd *cobra.Command, tag string) error {
	if !luxtorpedaTag.MatchString(tag) {
		return fmt.Errorf("No valid Luxtorpeda version tag: %q", tag)
	}

	version := tag[1:]

	var candidates []toolCandidate
	for _, t := range []string{tag, version} {
		name := "luxtorpeda-" + t
		candidates = append(candidates, toolCandidate{name: name, archive: name + ".tar.xz"})
	}

	return installTool(cmd.OutOrStdout(), toolRelease{
		candidates: candidates,
		baseURL:    githubURL + "/" + luxtorpedaRepo + "/releases/download/" + tag,
		strip:      true,
		postExtract: func(c toolCandidate, dir string) error {
			return patchLuxtorpedaVdf(dir, c.name, version)
		},
	})
}

// patchLuxtorpedaVdf sets the internal name of the tool to the directory name
// and appends the version to the display name, so that several versions can
// be installed side by side
func patchLuxtorpedaVdf(dir, name, version string) error {
	vdfPath := filepath.Join(dir, "compatibilitytool.vdf")
	root, err := steam.ParseTextConfig(vdfPath)
	if err != nil {
		return err
	}

	n := root.FirstByName("compat_tools").FirstSubTree()
	if n == nil {
		return fmt.Errorf("%v: no compatibility tool found", vdfPath)
	}

	n.SetName(name)

	displayName := n.FirstByName("display_name")
	if displayName != nil {
		displayName.SetString(displayName.String() + " " + version)
	}

	b, err := root.MarshalText()
	if err != nil {
		return err
	}

	return os.WriteFile(vdfPath, b, 0600)
}
