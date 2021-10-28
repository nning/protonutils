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

type mapLevel = map[string]interface{}

func getUID(u string) (string, error) {
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

		fmt.Fprintln(os.Stderr,
			"Warning: Several Steam users available, using "+uid+"\n"+
				"All available users: "+strings.Join(users, ", ")+"\n"+
				"Option \"-u\" can be used to specify user\n")
	}

	return uid, nil
}

func lookup(m mapLevel, x []string) (mapLevel, error) {
	y := m

	for _, s := range x {
		if y[s] == nil {
			return nil, errors.New("Key not found: " + s)
		}

		y = y[s].(mapLevel)
	}

	return y, nil
}

func vdfLookup(file string, x ...string) (mapLevel, error) {
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

func (s *Steam) cachedVdfLookup(cacheKey, file string, x ...string) (mapLevel, error) {
	m := s.vdfCache[cacheKey]
	if m != nil {
		return m.(mapLevel), nil
	}

	m, err := vdfLookup(file, x...)
	if err != nil {
		return nil, err
	}

	s.vdfCache[cacheKey] = m
	return m.(mapLevel), nil
}

func (s *Steam) getCompatToolMapping() (mapLevel, error) {
	return s.cachedVdfLookup("compatToolMapping", "config/config.vdf", "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
}

func (s *Steam) getLibraryConfig() (mapLevel, error) {
	return s.cachedVdfLookup("libraryConfig", "steamapps/libraryfolders.vdf", "libraryfolders")
}

func (s *Steam) getLocalConfig(user string) (mapLevel, error) {
	uid, err := getUID(user)
	if err != nil {
		return nil, err
	}

	return s.cachedVdfLookup("localConfig"+uid, "userdata/"+uid+"/config/localconfig.vdf", "UserLocalConfigStore", "Software", "Valve", "Steam", "apps")
}

func (s *Steam) getLoginUsers() (mapLevel, error) {
	return s.cachedVdfLookup("loginUsers", "config/loginusers.vdf", "users")
}

func (s *Steam) isInstalled(id string) (bool, error) {
	m, err := s.getLibraryConfig()
	if err != nil {
		return false, err
	}

	for i := 0; i < 10; i++ {
		x := m[fmt.Sprint(i)]
		if x == nil {
			break
		}

		apps := x.(mapLevel)["apps"].(mapLevel)
		for app := range apps {
			if app == id {
				return true, nil
			}
		}
	}

	return false, nil
}
