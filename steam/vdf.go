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

		fmt.Fprintln(os.Stderr, "Warning: Several Steam users available, using "+uid)
		fmt.Fprintln(os.Stderr, "All available users: "+strings.Join(users, ", ")+"\n")
		fmt.Fprintln(os.Stderr, "Option \"-u\" can be used to specify user")
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

func (s *Steam) cachedVdfLookup(cacheKey, file string, x ...string) (MapLevel, error) {
	m := s.vdfCache[cacheKey]

	if m == nil {
		m, err := vdfLookup(file, x...)
		if err != nil {
			return nil, err
		}
		s.vdfCache[cacheKey] = m
		return m, nil
	} else {
		return m.(MapLevel), nil
	}
}

func (s *Steam) GetCompatToolMapping() (MapLevel, error) {
	return s.cachedVdfLookup("compatToolMapping", "config/config.vdf", "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
}

func (s *Steam) GetLibraryConfig() (MapLevel, error) {
	return s.cachedVdfLookup("libraryConfig", "steamapps/libraryfolders.vdf", "libraryfolders")
}

func (s *Steam) GetLocalConfig(user string) (MapLevel, error) {
	uid, err := getUid(user)
	if err != nil {
		return nil, err
	}

	return s.cachedVdfLookup("localConfig"+uid, "userdata/"+uid+"/config/localconfig.vdf", "UserLocalConfigStore", "Software", "Valve", "Steam", "apps")
}

func (s *Steam) GetLoginUsers() (MapLevel, error) {
	return s.cachedVdfLookup("loginUsers", "config/loginusers.vdf", "users")
}

func (s *Steam) IsInstalled(id string) (bool, error) {
	m, err := s.GetLibraryConfig()
	if err != nil {
		return false, err
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
