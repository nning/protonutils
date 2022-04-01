package utils

import (
	"os"
	"path/filepath"
)

// DirSize retusn total size of directory
func DirSize(path string) (uint64, error) {
	var size uint64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += uint64(info.Size())
		}

		return err
	})

	return size, err
}
