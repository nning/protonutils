package steam

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/andygrunwald/vdf"
)

type mapLevel = map[string]interface{}

// KeyNotFoundError is returned if key in VDF is not found
type KeyNotFoundError struct {
	Name string
}

// GameInfo contains ID, Name, and LibraryPath of a game
type GameInfo struct {
	ID          string
	Name        string
	LibraryPath string
}

func (e *KeyNotFoundError) Error() string {
	return "Key not found: " + e.Name
}

func lookup(m mapLevel, x []string) (mapLevel, error) {
	y := m

	for _, s := range x {
		if y[s] == nil {
			return nil, &KeyNotFoundError{s}
		}

		y = y[s].(mapLevel)
	}

	return y, nil
}

func (s *Steam) vdfLookup(file string, x ...string) (mapLevel, error) {
	file = path.Join(s.Root, file)
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

	m, err := s.vdfLookup(file, x...)
	if err != nil {
		return nil, err
	}

	s.vdfCache[cacheKey] = m
	return m.(mapLevel), nil
}

func (s *Steam) getCompatToolMapping() (mapLevel, error) {
	key := []string{"InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping"}
	m, err := s.cachedVdfLookup("compatToolMapping", "config/config.vdf", key...)

	_, isKeyNotFoundError := err.(*KeyNotFoundError)
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
	return s.cachedVdfLookup("localConfig"+s.UID, "userdata/"+s.UID+"/config/localconfig.vdf", "UserLocalConfigStore", "Software", "Valve", "Steam", "apps")
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
func (s *Steam) GetGameInfo(idOrName string) (*GameInfo, error) {
	p, err := s.GetLibraryPathByID(idOrName)
	if err != nil {
		return nil, err
	}

	err = s.ReadCompatToolVersions()
	if err != nil {
		return nil, err
	}

	info := &GameInfo{idOrName, "", p}

	if p != "" {
		info.Name, _ = s.AppidCache.Get(info.ID)
		return info, nil
	}

	for _, version := range s.CompatToolVersions {
		for name, game := range version.Games {
			a := strings.ToLower(name)
			b := strings.ToLower(idOrName)

			if a == b || strings.HasPrefix(a, b) && game.IsInstalled {
				info = &GameInfo{game.ID, name, ""}
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
