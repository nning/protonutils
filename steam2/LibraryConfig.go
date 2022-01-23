package steam2

import (
	"path"
)

// LocalConfigVdf represents parsed VDF config for library config from
// libraryfolders.vdf
type LibraryConfigVdf struct {
	Vdf
}

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

func (vdf LibraryConfigVdf) IsInstalled(id string) bool {
	return vdf.GetLibraryPathByID(id) != ""
}

// GetLibraryConfig reads and parses libraryfolders.vdf and returns a
// LibraryConfigVdf
func (s *Steam) InitLibraryConfig() error {
	p := path.Join(s.Root, "steamapps", "libraryfolders.vdf")

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	s.LibraryConfig = &LibraryConfigVdf{Vdf{n, nil, p, s}}

	return nil
}
