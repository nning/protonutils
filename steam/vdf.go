package steam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/andygrunwald/vdf"

	. "github.com/nning/list_proton_versions"
)

type MapLevel = map[string]interface{}

func getUid(u string) string {
	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, ".steam", "root", "userdata")

	entries, err := ioutil.ReadDir(dir)
	PanicOnError(err)

	uid := entries[0].Name()

	if len(entries) > 1 {
		users := ""
		for i := 0; i < len(entries); i++ {
			name := entries[i].Name()
			if name == u {
				return name
			}

			comma := ", "
			if i == 0 {
				comma = ""
			}
			users = users + comma + name
		}

		fmt.Fprintln(os.Stderr, "Warning: Several Steam users available and only one is currently supported, using "+uid)
		fmt.Fprintln(os.Stderr, "All available users: "+users+"\n")
	}

	return uid
}

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
	file = path.Join(usr.HomeDir, ".steam", "root", file)

	f, err := os.Open(file)
	PanicOnError(err)

	p := vdf.NewParser(f)
	m, err := p.Parse()
	PanicOnError(err)

	return lookup(m, x)
}

func (s *Steam) GetCompatToolMapping() (MapLevel, error) {
	return vdfLookup("config/config.vdf", "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
}

func (s *Steam) GetLibraryConfig() (MapLevel, error) {
	return vdfLookup("steamapps/libraryfolders.vdf", "libraryfolders")
}

func (s *Steam) GetLocalConfig(user string) (MapLevel, error) {
	return vdfLookup("userdata/"+getUid(user)+"/config/localconfig.vdf", "UserLocalConfigStore", "Software", "Valve", "Steam", "apps")
}

func (s *Steam) IsInstalled(id string) bool {
	m := s.libraryConfig
	var err error

	if m == nil {
		m, err = s.GetLibraryConfig()
		PanicOnError(err)
		s.libraryConfig = m
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
