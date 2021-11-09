package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

// Feed represents Atom root node
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

// Entry represents Atom feed entry
type Entry struct {
	Link Link `xml:"link"`
}

// Link represents Atom feed entry link
type Link struct {
	URL string `xml:"href,attr"`
}

// WriteCounter implements printing download progress
type WriteCounter struct {
	Total    uint64
	Filename string
}

var egrollCmd = &cobra.Command{
	Use:   "ge",
	Short: "Commands for Proton-GE",
}

var egrollUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Download and extract the latest Proton-GE release",
	Run:   egrollUpdate,
}

var force bool
var keep bool

// Write counts bytes already written to wc
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints human-readable count of bytes already written to wc
func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading %s... %s complete", wc.Filename, humanize.Bytes(wc.Total))
}

func init() {
	rootCmd.AddCommand(egrollCmd)
	egrollCmd.AddCommand(egrollUpdateCmd)

	egrollUpdateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force last version update")
	egrollUpdateCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")
}

func getURL(url string) (io.Reader, error) {
	c := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func egrollUpdate(cmd *cobra.Command, args []string) {
	feedURL := "https://github.com/GloriousEggroll/proton-ge-custom/releases.atom"
	r, err := getURL(feedURL)
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
	dirpath := "Proton-" + tag
	filepath := dirpath + ".tar.gz"
	downloadURL := "https://github.com/GloriousEggroll/proton-ge-custom/releases/download/" + tag + "/" + filepath
	r, err = getURL(downloadURL)
	exitOnError(err)

	home, err := os.UserHomeDir()
	exitOnError(err)

	dir := path.Join(home, ".steam", "root", "compatibilitytools.d")
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

	out, err := os.Create(filepath)
	exitOnError(err)
	defer out.Close()

	counter := &WriteCounter{}
	counter.Filename = dirpath
	_, err = io.Copy(out, io.TeeReader(r, counter))
	exitOnError(err)

	fmt.Println("\nExtracting...")
	c := "tar xf " + filepath
	_, err = exec.Command("sh", "-c", c).Output()
	exitOnError(err)

	d, err := os.Getwd()
	fmt.Println(d, keep)

	if !keep {
		err = os.Remove(filepath)
		exitOnError(err)
	}
}
