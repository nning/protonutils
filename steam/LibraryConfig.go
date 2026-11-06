package steam

import (
	"log/slog"
	"path/filepath"
)

// LibraryConfigVdf represents parsed VDF config for library config from
// libraryfolders.vdf
type LibraryConfigVdf struct {
	Vdf
}

// GetLibraryPaths returns array of Steam library paths
func (v LibraryConfigVdf) GetLibraryPaths() []string {
	var paths []string

	for x := v.Root.FirstSubTree(); x != nil; x = x.NextSubTree() {
		if p := x.FirstByName("path").String(); p != "" {
			paths = append(paths, p)
		}
	}

	return paths
}

// GetLibraryPathByID returns library path for app id (empty string if app is
// not installed)
func (v LibraryConfigVdf) GetLibraryPathByID(id string) string {
	if id == "" {
		return ""
	}

	for x := v.Root.FirstSubTree(); x != nil; x = x.NextSubTree() {
		if x.FirstByName("apps").FirstByName(id) != nil {
			return x.FirstByName("path").String()
		}
	}

	return ""
}

// IsInstalled returns whether app id is installed
func (v LibraryConfigVdf) IsInstalled(id string) bool {
	return v.GetLibraryPathByID(id) != ""
}

func (s *Steam) initLibraryConfig() error {
	p := filepath.Join(s.Root, "steamapps", "libraryfolders.vdf")
	slog.Debug("steam.initLibraryConfig", "path", p)

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	s.LibraryConfig = &LibraryConfigVdf{Vdf{n, nil, p, s}}

	return nil
}
