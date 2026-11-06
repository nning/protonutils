package steam

import (
	"errors"
	"path/filepath"
)

// GetCompatdataPath returns compatdata path and game name for given game ID or
// name. For games that are not installed (anymore) and non-Steam shortcuts, the
// library folders are searched for an existing compatdata directory.
func (s *Steam) GetCompatdataPath(idOrName string) (string, string, error) {
	id, name, err := s.GetAppIDAndName(idOrName)
	if err != nil {
		return "", "", err
	}

	if p := s.LibraryConfig.GetLibraryPathByID(id); p != "" {
		return filepath.Join(p, "steamapps", "compatdata", id), name, nil
	}

	if p := s.SearchCompatdataPath(id); p != "" {
		return p, name, nil
	}

	return "", "", errors.New("Game not installed and no compatdata directory found")
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
