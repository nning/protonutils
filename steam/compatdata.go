package steam

import (
	"errors"
	"os"
	"path"
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

	return path.Join(p, "steamapps", "compatdata", id), name, nil
}

// SearchCompatdataPath searches for compatdata path for given game id in all
// library paths
func (s *Steam) SearchCompatdataPath(id string) string {
	paths := s.LibraryConfig.GetLibraryPaths()
	for _, p := range paths {
		x := path.Join(p, "steamapps", "compatdata", id)

		info, err := os.Stat(x)
		if err != nil {
			continue
		}

		if info.IsDir() {
			return x
		}
	}

	return ""
}
