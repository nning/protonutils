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
	"strconv"
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
	Finished uint64
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
	wc.Finished += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints human-readable count of bytes already written to wc
func (wc *WriteCounter) PrintProgress() {
	p := uint64(float64(wc.Finished) / float64(wc.Total) * 100)

	fmt.Printf("\r%s", strings.Repeat(" ", 80))
	fmt.Printf("\rDownloading %s... %d%% (%s of %s) complete", wc.Filename, p, humanize.Bytes(wc.Finished), humanize.Bytes(wc.Total))
}

func init() {
	rootCmd.AddCommand(egrollCmd)
	egrollCmd.AddCommand(egrollUpdateCmd)

	egrollUpdateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force last version update")
	egrollUpdateCmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep downloaded archive of last version")
}

func getURL(url string) (io.Reader, uint64, error) {
	c := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}

	size, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	return res.Body, uint64(size), nil
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
	dirpath := "Proton-" + tag
	filepath := dirpath + ".tar.gz"

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

	d, err := os.Getwd()
	fmt.Println(d, keep)

	if !keep {
		err = os.Remove(filepath)
		exitOnError(err)
	}
}
