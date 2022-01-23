package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/nning/protonutils/steam"
	"github.com/nning/protonutils/steam2"
	"github.com/spf13/cobra"
)

var luxtorpedaCmd = &cobra.Command{
	Use:   "luxtorpeda",
	Short: "Commands for Luxtorpeda",
}

var luxtorpedaDownloadCmd = &cobra.Command{
	Use:   "download [flags] <version>",
	Short: "Download and extract version specified in argument",
	Args:  cobra.MinimumNArgs(1),
	Run:   luxtorpedaDownload,
}

var luxtorpedaUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Download and extract the latest Luxtorpeda release",
	Run:   luxtorpedaUpdate,
}

func init() {
	rootCmd.AddCommand(luxtorpedaCmd)
	luxtorpedaCmd.AddCommand(luxtorpedaDownloadCmd)
	luxtorpedaCmd.AddCommand(luxtorpedaUpdateCmd)

	luxtorpedaDownloadCmd.Flags().BoolVarP(&extractOnly, "extract-only", "e", false, "Do not download but extract only from existing archive")
	luxtorpedaDownloadCmd.Flags().BoolVarP(&force, "force", "f", false, "Force download even if version exists")
	luxtorpedaDownloadCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")

	luxtorpedaUpdateCmd.Flags().BoolVarP(&extractOnly, "extract-only", "e", false, "Do not download but extract only from existing archive")
	luxtorpedaUpdateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force last version update")
	luxtorpedaUpdateCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")
}

func luxtorpedaDownload(cmd *cobra.Command, args []string) {
	tag := args[0]

	valid, err := regexp.MatchString("^v[0-9]*", tag)
	exitOnError(err)

	if !valid {
		fmt.Fprintln(os.Stderr, "No valid Luxtorpeda version tag")
		os.Exit(1)
	}

	dirpath := "luxtorpeda-" + tag[1:]
	filepath := dirpath + ".tar.xz"

	s, err := steam.New(user, cfg.SteamRoot, false)
	exitOnError(err)

	dir := getCompatDir(s)
	_, err = os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, 0700)
		exitOnError(err)
	}

	err = os.Chdir(dir)
	exitOnError(err)

	stat, err := os.Stat(dirpath)
	if err == nil && stat.IsDir() && !force {
		fmt.Println(dirpath, "already available")
		return
	}

	if force {
		err := os.RemoveAll(dirpath)
		exitOnError(err)
	}

	if extractOnly {
		_, err := os.Stat(filepath)
		exitOnError(err)
	} else {
		downloadURL := "https://github.com/luxtorpeda-dev/luxtorpeda/releases/download/" + tag + "/" + filepath
		r, size, err := getURL(downloadURL)
		exitOnError(err)

		out, err := os.Create(filepath)
		exitOnError(err)
		defer out.Close()

		counter := &WriteCounter{}
		counter.Filename = dirpath
		counter.Total = size
		_, err = io.Copy(out, io.TeeReader(r, counter))
		exitOnError(err)
	}

	fmt.Println("\nExtracting...")
	err = os.Mkdir(dirpath, 0700)
	exitOnError(err)
	c := "tar xf " + filepath + " -C " + dirpath + " --strip-components 1"
	_, err = exec.Command("sh", "-c", c).Output()
	exitOnError(err)

	vdfPath := path.Join(dirpath, "compatibilitytool.vdf")
	root, err := steam2.ParseTextConfig(vdfPath)
	exitOnError(err)

	n := root.FirstChild().FirstChild()
	n.SetName(dirpath)

	n = n.FirstByName("display_name")
	v := n.String() + " " + tag[1:]
	n.SetString(v)

	b, err := root.MarshalText()
	exitOnError(err)

	err = os.WriteFile(vdfPath, b, 0600)
	exitOnError(err)

	if !keep {
		err = os.Remove(filepath)
		exitOnError(err)
	}
}

func luxtorpedaUpdate(cmd *cobra.Command, args []string) {
	feedURL := "https://github.com/luxtorpeda-dev/luxtorpeda/releases.atom"
	r, _, err := getURL(feedURL)
	exitOnError(err)

	body, err := ioutil.ReadAll(r)
	exitOnError(err)

	var feed Feed
	err = xml.Unmarshal(body, &feed)
	exitOnError(err)

	if len(feed.Entries) == 0 {
		exitOnError(nil, "Could not fetch releases")
	}

	releaseURL := feed.Entries[0].Link.URL
	if releaseURL == "" {
		exitOnError(nil, "Could not get URL for latest release")
	}

	tag := path.Base(releaseURL)
	luxtorpedaDownload(cmd, []string{tag})
}
