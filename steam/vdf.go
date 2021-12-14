package steam

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/andygrunwald/vdf"
)

type mapLevel = map[string]interface{}

type keyNotFoundError struct {
	name string
}

type gameInfo struct {
	ID          string
	Name        string
	LibraryPath string
}

func (e *keyNotFoundError) Error() string {
	return "Key not found: " + e.name
}

func lookup(m mapLevel, x []string) (mapLevel, error) {
	y := m

	for _, s := range x {
		if y[s] == nil {
			return nil, &keyNotFoundError{s}
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
	key := []string{"InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping"}
	m, err := s.cachedVdfLookup("compatToolMapping", "config/config.vdf", key...)

	_, isKeyNotFoundError := err.(*keyNotFoundError)
	if err != nil && isKeyNotFoundError {
		key[3] = "steam"
		m, err = s.cachedVdfLookup("compatToolMapping", "config/config.vdf", key...)
	}

	return m, err
}

func (s *Steam) getLibraryConfig() (mapLevel, error) {
	return s.cachedVdfLookup("libraryConfig", "steamapps/libraryfolders.vdf", "libraryfolders")
}

func (s *Steam) getLocalConfig() (mapLevel, error) {
	return s.cachedVdfLookup("localConfig"+s.uid, "userdata/"+s.uid+"/config/localconfig.vdf", "UserLocalConfigStore", "Software", "Valve", "Steam", "apps")
}

func (s *Steam) getLoginUsers() (mapLevel, error) {
	return s.cachedVdfLookup("loginUsers", "config/loginusers.vdf", "users")
}

func (s *Steam) isInstalled(id string) (bool, error) {
	path, err := s.GetLibraryPathByID(id)
	if err != nil || path == "" {
		return false, err
	}

	return true, nil
}

// GetLibraryPathByID returns path to Steam library folder, the game with specified app ID is installed in
func (s *Steam) GetLibraryPathByID(id string) (string, error) {
	m, err := s.getLibraryConfig()
	if err != nil {
		return "", err
	}

	for _, x := range m {
		v, isMapLevel := x.(mapLevel)
		if !isMapLevel {
			continue
		}

		for app := range v["apps"].(mapLevel) {
			if app == id {
				return x.(mapLevel)["path"].(string), nil
			}
		}
	}

	return "", nil
}

// GetGameInfo returns library path (and ID and name) by game ID or name
func (s *Steam) GetGameInfo(idOrName string) (*gameInfo, error) {
	p, err := s.GetLibraryPathByID(idOrName)
	if err != nil {
		return nil, err
	}

	err = s.ReadCompatToolVersions()
	if err != nil {
		return nil, err
	}

	info := &gameInfo{idOrName, "", p}

	if p != "" {
		info.Name, _ = s.AppidCache.Get(info.ID)
		return info, nil
	}

	for _, games := range s.CompatToolVersions {
		for name, game := range games {
			a := strings.ToLower(name)
			b := strings.ToLower(idOrName)

			if a == b || strings.HasPrefix(a, b) && game.IsInstalled {
				info = &gameInfo{game.ID, name, ""}
				break
			}
		}
	}

	p, err = s.GetLibraryPathByID(info.ID)
	if err != nil {
		return nil, err
	}

	info.LibraryPath = p

	if info.ID == "" || p == "" {
		fmt.Fprintln(os.Stderr, "App ID or path not found")
		os.Exit(1)
	}

	return info, nil
}
