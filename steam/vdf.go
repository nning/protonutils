package steam

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/andygrunwald/vdf"

	. "github.com/nning/list_proton_versions"
)

type MapLevel = map[string]interface{}

func lookup(m MapLevel, x []string) (MapLevel, error) {
	y := m

	for _, s := range x {
		if y[s] == nil {
			return nil, errors.New("Key not found: " + s)
		} else {
			y = y[s].(MapLevel)
		}
	}

	return y, nil
}

func vdfLookup(file string, x ...string) (MapLevel, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	file = path.Join(dir, ".steam", "root", file)

	f, err := os.Open(file)
	PanicOnError(err)

	p := vdf.NewParser(f)
	m, err := p.Parse()
	PanicOnError(err)

	return lookup(m, x)
}

func (self *Steam) GetCompatToolMapping(x ...string) (MapLevel, error) {
	return vdfLookup("config/config.vdf", "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
}

func (self *Steam) GetLibraryConfig() (MapLevel, error) {
	return vdfLookup("steamapps/libraryfolders.vdf", "libraryfolders")
}

func (self *Steam) IsInstalled(id string) bool {
	m := self.libraryConfig
	var err error

	if m == nil {
		m, err = self.GetLibraryConfig()
		PanicOnError(err)
		self.libraryConfig = m
	}

	for i := 0; i < 10; i++ {
		x := m[fmt.Sprint(i)]
		if x == nil {
			break
		}

		apps := x.(MapLevel)["apps"].(MapLevel)
		for app := range apps {
			if app == id {
				return true
			}
		}
	}

	return false
}
