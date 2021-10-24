package steam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/andygrunwald/vdf"
)

type MapLevel = map[string]interface{}

func getUid(u string) (string, error) {
	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, ".steam", "root", "userdata")

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	uid := entries[0].Name()

	if len(entries) > 1 {
		users := make([]string, len(entries))

		for i, entry := range entries {
			name := entry.Name()
			if name == u {
				return name, nil
			}

			users[i] = name
		}

		fmt.Fprintln(os.Stderr, "Warning: Several Steam users available and only one is currently supported, using "+uid)
		fmt.Fprintln(os.Stderr, "All available users: "+strings.Join(users, ", ")+"\n")
	}

	return uid, nil
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
	if err != nil {
		return nil, err
	}

	p := vdf.NewParser(f)
	m, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return lookup(m, x)
}

func (s *Steam) GetCompatToolMapping() (MapLevel, error) {
	return vdfLookup("config/config.vdf", "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
}

func (s *Steam) GetLibraryConfig() (MapLevel, error) {
	return vdfLookup("steamapps/libraryfolders.vdf", "libraryfolders")
}

func (s *Steam) GetLocalConfig(user string) (MapLevel, error) {
	uid, err := getUid(user)
	if err != nil {
		return nil, err
	}

	return vdfLookup("userdata/"+uid+"/config/localconfig.vdf", "UserLocalConfigStore", "Software", "Valve", "Steam", "apps")
}

func (s *Steam) GetLoginUsers() (MapLevel, error) {
	return vdfLookup("config/loginusers.vdf", "users")
}

func (s *Steam) IsInstalled(id string) (bool, error) {
	m := s.libraryConfig
	var err error

	if m == nil {
		m, err = s.GetLibraryConfig()
		if err != nil {
			return false, err
		}
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
				return true, nil
			}
		}
	}

	return false, nil
}
