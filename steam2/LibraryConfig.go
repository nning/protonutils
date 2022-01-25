package steam2

import (
	"path"
)

// LibraryConfigVdf represents parsed VDF config for library config from
// libraryfolders.vdf
type LibraryConfigVdf struct {
	Vdf
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

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	s.LibraryConfig = &LibraryConfigVdf{Vdf{n, nil, p, s}}

	return nil
}
