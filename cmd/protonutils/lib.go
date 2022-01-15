package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/nning/protonutils/steam"
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

	if res.StatusCode != 200 {
		return nil, 0, errors.New("Error retrieving URL (" + res.Status + "):\n" + url)
	}

	size, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	return res.Body, uint64(size), nil
}

func getCompatDir(s *steam.Steam) string {
	return path.Join(s.Root, "compatibilitytools.d")
}
