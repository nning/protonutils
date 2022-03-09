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
	"github.com/spf13/cobra"
)

var egrollCmd = &cobra.Command{
	Use:   "ge",
	Short: "Commands for Proton-GE",
}

var egrollDownloadCmd = &cobra.Command{
	Use:   "download [flags] <version>",
	Short: "Download and extract version specified in argument",
	Args:  cobra.MinimumNArgs(1),
	Run:   egrollDownload,
}

var egrollUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Download and extract the latest Proton-GE release",
	Run:   egrollUpdate,
}

func init() {
	rootCmd.AddCommand(egrollCmd)
	egrollCmd.AddCommand(egrollDownloadCmd)
	egrollCmd.AddCommand(egrollUpdateCmd)

	egrollDownloadCmd.Flags().BoolVarP(&force, "force", "f", false, "Force download even if version exists")
	egrollDownloadCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")

	egrollUpdateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force last version update")
	egrollUpdateCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")
}

func egrollDownload(cmd *cobra.Command, args []string) {
	tag := args[0]

	validOld, err := regexp.MatchString("^[0-9]*\\.[0-9]*(rc[0-9])?-GE-[0-9]*", tag)
	validNew, err := regexp.MatchString("^GE-Proton[0-9]*-[0-9]*", tag)
	exitOnError(err)

	if !validOld && !validNew {
		fmt.Fprintln(os.Stderr, "No valid GE version tag")
		os.Exit(1)
	}

	dirpath := tag
	if validOld {
		dirpath = "Proton-" + tag
	}
	filepath := dirpath + ".tar.gz"

	s, err := steam.New(user, cfg.SteamRoot, false)
	exitOnError(err)

	dir := s.GetCompatibilityToolsDir()
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

	downloadURL := "https://github.com/GloriousEggroll/proton-ge-custom/releases/download/" + tag + "/" + filepath
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

	fmt.Println("\nExtracting...")
	c := "tar xf " + filepath
	_, err = exec.Command("sh", "-c", c).Output()
	exitOnError(err)

	if !keep {
		err = os.Remove(filepath)
		exitOnError(err)
	}
}

func egrollUpdate(cmd *cobra.Command, args []string) {
	feedURL := "https://github.com/GloriousEggroll/proton-ge-custom/releases.atom"
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
	egrollDownload(cmd, []string{tag})
}
