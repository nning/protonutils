package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// githubURL is the base URL for release downloads (variable for tests)
var githubURL = "https://github.com"

// tmpDirPrefix is the prefix of temporary directories in compatibilitytools.d
const tmpDirPrefix = ".protonutils-tmp-"

var downloadClient = &http.Client{
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	},
}

// errNotFound is returned by httpGet if the server responds with 404
var errNotFound = errors.New("not found")

// feed represents Atom root node
type feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []struct {
		Link struct {
			URL string `xml:"href,attr"`
		} `xml:"link"`
	} `xml:"entry"`
}

// progressWriter prints download progress (at most once per percent or MB)
type progressWriter struct {
	w        io.Writer
	name     string
	total    uint64
	finished uint64
	last     uint64
}

func (p *progressWriter) Write(b []byte) (int, error) {
	p.finished += uint64(len(b))

	step := p.finished / (1 << 20)
	if p.total > 0 {
		step = p.finished * 100 / p.total
	}

	if step != p.last || p.finished == p.total {
		p.last = step
		p.print()
	}

	return len(b), nil
}

func (p *progressWriter) print() {
	if p.total > 0 {
		fmt.Fprintf(p.w, "\r\033[KDownloading %s... %d%% (%s of %s) complete", p.name, p.finished*100/p.total, humanize.Bytes(p.finished), humanize.Bytes(p.total))
	} else {
		fmt.Fprintf(p.w, "\r\033[KDownloading %s... %s complete", p.name, humanize.Bytes(p.finished))
	}
}

func httpGet(url string) (*http.Response, error) {
	res, err := downloadClient.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		res.Body.Close()
		return nil, fmt.Errorf("%w: %v", errNotFound, url)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("Error retrieving URL (%v):\n%v", res.Status, url)
	}

	return res, nil
}

func httpGetBytes(url string) ([]byte, error) {
	res, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(io.LimitReader(res.Body, 10<<20))
}

// downloadFile downloads url to file dest and returns the SHA-512 checksum of
// the content
func downloadFile(w io.Writer, url, dest, name string) ([]byte, error) {
	res, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return nil, err
	}

	var total uint64
	if res.ContentLength > 0 {
		total = uint64(res.ContentLength)
	}

	h := sha512.New()
	progress := &progressWriter{w: w, name: name, total: total}
	_, err = io.Copy(io.MultiWriter(f, h, progress), res.Body)
	fmt.Fprintln(w)

	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// latestReleaseTag returns the tag of the latest release of a GitHub repository
func latestReleaseTag(repo string) (string, error) {
	body, err := httpGetBytes(githubURL + "/" + repo + "/releases.atom")
	if err != nil {
		return "", err
	}

	var f feed
	err = xml.Unmarshal(body, &f)
	if err != nil {
		return "", err
	}

	if len(f.Entries) == 0 {
		return "", errors.New("Could not fetch releases")
	}

	releaseURL := f.Entries[0].Link.URL
	if releaseURL == "" {
		return "", errors.New("Could not get URL for latest release")
	}

	return path.Base(releaseURL), nil
}

// verifyChecksum compares sum with the SHA-512 checksum in the file at
// checksumURL (sha512sum format). If no checksum file exists, a warning is
// printed.
func verifyChecksum(w io.Writer, checksumURL string, sum []byte) error {
	body, err := httpGetBytes(checksumURL)
	if errors.Is(err, errNotFound) {
		fmt.Fprintln(w, "WARNING: No checksum available, skipping verification")
		return nil
	}
	if err != nil {
		return err
	}

	fields := strings.Fields(string(body))
	if len(fields) == 0 {
		return errors.New("checksum file is empty")
	}

	expected, err := hex.DecodeString(fields[0])
	if err != nil || len(expected) != sha512.Size {
		return errors.New("checksum file is invalid")
	}

	if hex.EncodeToString(expected) != hex.EncodeToString(sum) {
		return errors.New("checksum mismatch, download is corrupt")
	}

	fmt.Fprintln(w, "Checksum OK")
	return nil
}

// extractTar extracts archive into dir (using the tar command, which detects
// the compression)
func extractTar(archive, dir string, stripComponents int) error {
	args := []string{"-x", "-f", archive, "-C", dir}
	if stripComponents > 0 {
		args = append(args, fmt.Sprintf("--strip-components=%d", stripComponents))
	}

	out, err := exec.Command("tar", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("extracting %v: %w: %s", filepath.Base(archive), err, strings.TrimSpace(string(out)))
	}

	return nil
}

// toolCandidate is a possible combination of directory name (in
// compatibilitytools.d) and archive file name of a release
type toolCandidate struct {
	name    string
	archive string
}

// toolRelease describes a compatibility tool release archive
type toolRelease struct {
	// candidates are tried in order (archive names relative to baseURL or
	// compatibilitytools.d)
	candidates []toolCandidate
	// baseURL for archive and checksum downloads
	baseURL string
	// checksum returns the checksum file name (relative to baseURL) for a
	// candidate (optional)
	checksum func(c toolCandidate) string
	// strip first path component and extract into directory name
	strip bool
	// postExtract is called with the candidate and the extracted directory
	// (optional)
	postExtract func(c toolCandidate, dir string) error
}

// installTool downloads and extracts a compatibility tool release into
// compatibilitytools.d. Everything is prepared in a temporary directory and
// only moved into place after successful download, verification and extraction.
func installTool(w io.Writer, r toolRelease) error {
	compatDir := filepath.Join(root, "compatibilitytools.d")

	for _, c := range r.candidates {
		if isDir(filepath.Join(compatDir, c.name)) && !force {
			fmt.Fprintln(w, c.name, "already available")
			return nil
		}
	}

	err := os.MkdirAll(compatDir, 0700)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(compatDir, tmpDirPrefix)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	var c toolCandidate
	var archive, keepPath string

	if extractOnly {
		for _, candidate := range r.candidates {
			p := filepath.Join(compatDir, candidate.archive)
			if _, err := os.Stat(p); err == nil {
				c, archive = candidate, p
				break
			}
		}

		if archive == "" {
			var names []string
			for _, candidate := range r.candidates {
				names = append(names, candidate.archive)
			}

			return fmt.Errorf("archive not found in %v (tried %v)", compatDir, strings.Join(names, ", "))
		}
	} else {
		var sum []byte

		for i, candidate := range r.candidates {
			c = candidate
			archive = filepath.Join(tmp, c.archive)
			sum, err = downloadFile(w, r.baseURL+"/"+c.archive, archive, c.name)
			if errors.Is(err, errNotFound) && i < len(r.candidates)-1 {
				continue
			}
			if err != nil {
				return err
			}

			keepPath = filepath.Join(compatDir, c.archive)
			break
		}

		if r.checksum != nil {
			err = verifyChecksum(w, r.baseURL+"/"+r.checksum(c), sum)
			if err != nil {
				return err
			}
		}
	}

	fmt.Fprintln(w, "Extracting...")

	extractDir := filepath.Join(tmp, "extract")
	extracted := filepath.Join(extractDir, c.name)
	strip := 0
	if r.strip {
		extractDir = extracted
		strip = 1
	}

	err = os.MkdirAll(extractDir, 0700)
	if err != nil {
		return err
	}

	err = extractTar(archive, extractDir, strip)
	if err != nil {
		return err
	}

	if !isDir(extracted) {
		return fmt.Errorf("unexpected archive content: %v not found", c.name)
	}

	if r.postExtract != nil {
		err = r.postExtract(c, extracted)
		if err != nil {
			return err
		}
	}

	target := filepath.Join(compatDir, c.name)
	old := filepath.Join(tmp, "old")
	hasOld := isDir(target)
	if hasOld {
		err = os.Rename(target, old)
		if err != nil {
			return err
		}
	}

	err = os.Rename(extracted, target)
	if err != nil {
		if hasOld {
			if restoreErr := os.Rename(old, target); restoreErr != nil {
				return errors.Join(err, fmt.Errorf("restoring previous version: %w", restoreErr))
			}
		}

		return err
	}

	if keep && keepPath != "" {
		err = os.Rename(archive, keepPath)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(w, "Installed", c.name)
	return nil
}

// isDir returns whether p is a directory (following symlinks)
func isDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}
