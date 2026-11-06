package steam

import (
	"errors"
	"path/filepath"
)

// GetCompatdataPath returns compatdata path and game name for given game ID or
// name
func (s *Steam) GetCompatdataPath(idOrName string) (string, string, error) {
	id, name, err := s.GetAppIDAndName(idOrName)
	if err != nil {
		return "", "", err
	}

	p := s.LibraryConfig.GetLibraryPathByID(id)
	if p == "" {
		return "", "", errors.New("Game not installed")
	}

	return filepath.Join(p, "steamapps", "compatdata", id), name, nil
}

// GetInstalldirPath returns installation path and game name for given game ID
// or name
func (s *Steam) GetInstalldirPath(idOrName string) (string, string, error) {
	id, name, err := s.GetAppIDAndName(idOrName)
	if err != nil {
		return "", "", err
	}

	p := s.LibraryConfig.GetLibraryPathByID(id)
	if p == "" {
		return "", "", errors.New("Game not installed")
	}

	dir, err := s.GetInstalldir(id)
	if err != nil {
		return "", "", err
	}

	if dir == "" {
		return "", "", errors.New("installdir not found")
	}

	return filepath.Join(p, "steamapps", "common", dir), name, nil
}

// SearchCompatdataPath searches for compatdata path for given game id in all
// library paths
func (s *Steam) SearchCompatdataPath(id string) string {
	for _, p := range s.LibraryConfig.GetLibraryPaths() {
		x := filepath.Join(p, "steamapps", "compatdata", id)
		if isDir(x) {
			return x
		}
	}

	return ""
}

// GetGamePaths returns the existing compatdata and installation directories of
// app id (empty strings if they do not exist)
func (s *Steam) GetGamePaths(id string) (compatdata, install string) {
	lib := s.LibraryConfig.GetLibraryPathByID(id)

	if lib != "" {
		p := filepath.Join(lib, "steamapps", "compatdata", id)
		if isDir(p) {
			compatdata = p
		}
	}

	if compatdata == "" {
		compatdata = s.SearchCompatdataPath(id)
	}

	if lib != "" && !IsShortcut(id) {
		if dir, err := s.GetInstalldir(id); err == nil && dir != "" {
			p := filepath.Join(lib, "steamapps", "common", dir)
			if isDir(p) {
				install = p
			}
		}
	}

	return compatdata, install
}
