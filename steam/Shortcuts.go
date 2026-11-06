package steam

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dustin/go-humanize"
)

// Shortcuts holds names of non-Steam game shortcuts (shortcuts.vdf) by app ID
type Shortcuts struct {
	Path  string
	names map[uint32]string
}

// ParseShortcuts parses the content of shortcuts.vdf
func ParseShortcuts(data []byte, path string) (*Shortcuts, error) {
	sc := &Shortcuts{Path: path, names: make(map[uint32]string)}
	if len(data) == 0 {
		return sc, nil
	}

	n, err := ParseBinaryVdf(data)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", path, err)
	}

	for x := n.FirstSubTree(); x != nil; x = x.NextSubTree() {
		id := x.FirstByName("appid")
		if id == nil {
			continue
		}

		sc.names[uint32(id.Int())] = x.FirstByName("AppName").String()
	}

	return sc, nil
}

// GetName returns the name of a shortcut by app ID (empty string if not found)
func (sc *Shortcuts) GetName(id string) string {
	appID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ""
	}

	return sc.names[uint32(appID)]
}

func (s *Steam) initShortcuts() error {
	p := filepath.Join(s.Root, "userdata", s.UID, "config", "shortcuts.vdf")
	in, err := os.ReadFile(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	slog.Debug("Steam.initShortcuts", "path", p, "size", humanize.Bytes(uint64(len(in))))

	s.Shortcuts, err = ParseShortcuts(in, p)
	if err != nil {
		fmt.Fprintf(s.Stderr, "Warning: %v, shortcut names are missing\n\n", err)
		s.Shortcuts, _ = ParseShortcuts(nil, p)
	}

	return nil
}
