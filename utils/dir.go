package utils

import (
	"io/fs"
	"path/filepath"
)

// DirSize returns total size of all regular files in directory
func DirSize(path string) (uint64, error) {
	var size uint64

	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			size += uint64(info.Size())
		}

		return nil
	})

	return size, err
}
