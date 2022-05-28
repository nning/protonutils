package steam

import (
	"path"

	log "github.com/sirupsen/logrus"
)

// LibraryConfigVdf represents parsed VDF config for library config from
// libraryfolders.vdf
type LibraryConfigVdf struct {
	Vdf
}

// GetLibraryPaths returns array of Steam library paths
func (vdf LibraryConfigVdf) GetLibraryPaths() []string {
	var paths []string

	x := vdf.Root.FirstSubTree()
	for {
		paths = append(paths, x.FirstByName("path").String())

		x = x.NextSubTree()
		if x == nil {
			break
		}
	}

	return paths
}

// GetLibraryPathByID returns library path for app id
func (vdf LibraryConfigVdf) GetLibraryPathByID(id string) string {
	x := vdf.Root.FirstSubTree()
	for {
		z := x.FirstByName("apps").FirstChild()
		for {
			if z.Name() == id {
				return x.FirstByName("path").String()
			}

			z = z.NextChild()
			if z == nil {
				break
			}
		}

		x = x.NextSubTree()
		if x == nil {
			break
		}
	}

	return ""
}

// IsInstalled returns whether app id is installed
func (vdf LibraryConfigVdf) IsInstalled(id string) bool {
	return vdf.GetLibraryPathByID(id) != ""
}

func (s *Steam) initLibraryConfig() error {
	p := path.Join(s.Root, "steamapps", "libraryfolders.vdf")
	log.Debug("steam.initLibraryConfig(", p, ")")

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	s.LibraryConfig = &LibraryConfigVdf{Vdf{n, nil, p, s}}

	return nil
}
