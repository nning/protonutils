package steam

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BenLubar/vdf"
	"github.com/dustin/go-humanize"
)

// BackupSuffix is appended to the file name of backups created by Vdf.Save
const BackupSuffix = ".protonutils.bak"

// KeyNotFoundError is returned if key in VDF is not found
type KeyNotFoundError struct {
	Name string
}

func (e *KeyNotFoundError) Error() string {
	return "Key not found: " + e.Name
}

// Vdf represents a parsed VDF file
type Vdf struct {
	Root  *vdf.Node
	Node  *vdf.Node
	Path  string
	Steam *Steam
}

// Save saves a parsed VDF file back to disk. The previous content is kept as
// backup (Path + BackupSuffix) and the file is replaced atomically.
func (v *Vdf) Save() error {
	out, err := v.Root.MarshalText()
	if err != nil {
		return err
	}

	var mode os.FileMode = 0600
	if info, err := os.Stat(v.Path); err == nil {
		mode = info.Mode().Perm()

		old, err := os.ReadFile(v.Path)
		if err != nil {
			return err
		}

		err = writeFileAtomic(v.Path+BackupSuffix, old, mode)
		if err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
	}

	return writeFileAtomic(v.Path, out, mode)
}

// writeFileAtomic writes data to a temporary file in the directory of p and
// renames it to p
func writeFileAtomic(p string, data []byte, mode os.FileMode) error {
	f, err := os.CreateTemp(filepath.Dir(p), "."+filepath.Base(p)+".tmp-*")
	if err != nil {
		return err
	}

	tmp := f.Name()
	defer os.Remove(tmp)

	_, err = f.Write(data)
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}

	err = os.Chmod(tmp, mode)
	if err != nil {
		return err
	}

	return os.Rename(tmp, p)
}

// Lookup looks up a "key path" in a parsed VDF tree
func Lookup(n *vdf.Node, x []string) (*vdf.Node, error) {
	y := n

	for _, key := range x {
		y = y.FirstByName(key)
		if y == nil {
			return nil, &KeyNotFoundError{Name: key}
		}
	}

	return y, nil
}

// ParseTextConfig reads a file and parses it as text VDF
func ParseTextConfig(p string) (*vdf.Node, error) {
	in, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	slog.Debug("steam.ParseTextConfig", "path", p, "size", humanize.Bytes(uint64(len(in))))

	var n vdf.Node
	err = n.UnmarshalText(in)
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", p, err)
	}

	return &n, nil
}
