package steam

import (
	"io"
	"os"
	"path"
	"sort"
)

// Version represents a compatibility tool version
type Version struct {
	ID    string
	Name  string
	Games Games
}

// CompatToolVersions maps Proton versions to games
type CompatToolVersions map[string]*Version

func (versions CompatToolVersions) includesGameID(id string) bool {
	for _, version := range versions {
		if version.Games.includesID(id) {
			return true
		}
	}

	return false
}

// Sort returns slice of alphabetically sorted CompatToolVersion names
func (versions CompatToolVersions) Sort() []string {
	keys := make([]string, len(versions))

	i := 0
	for key := range versions {
		keys[i] = key
		i++
	}

	sort.Strings(keys)
	return keys
}

func (s *Steam) includesGameID(id string) bool {
	return s.CompatToolVersions.includesGameID(id)
}

func (s *Steam) getCompatToolName(shortName string) string {
	if shortName == "" {
		return ""
	}

	str, _ := s.VersionNameCache.Get(shortName)
	if str != "" {
		return str
	}

	displayName, err := s.findCompatToolName(shortName)
	if err != nil || displayName == "" {
		s.VersionNameCache.Add(shortName, shortName, false)
		return shortName
	}

	s.VersionNameCache.Add(shortName, displayName, true)
	return displayName
}

// ReadCompatToolVersions reads Proton versions and games from different Steam configs
func (s *Steam) ReadCompatToolVersions() error {
	x, err := s.getCompatToolMapping()
	if err != nil {
		return err
	}

	defID := x["0"].(mapLevel)["name"].(string)
	def := s.getCompatToolName(defID) + " (Default)"

	for id, cfg := range x {
		vID := cfg.(mapLevel)["name"].(string)
		v := s.getCompatToolName(vID)
		if v == "" {
			v = def
		}

		_, err = s.addGame(vID, v, id)
		if err != nil && err != io.EOF {
			return err
		}
	}

	x, err = s.getLocalConfig()
	_, isKeyNotFoundError := err.(*keyNotFoundError)
	if err != nil && !isKeyNotFoundError {
		return err
	}

	for id, cfg := range x {
		v := cfg.(mapLevel)["ViewedSteamPlay"]
		if v == nil {
			continue
		}

		if !s.includesGameID(id) {
			_, err = s.addGame(defID, def, id)
			if err != nil && err != io.EOF {
				return err
			}
		}
	}

	return nil
}

func (s *Steam) GetGameVersion(id string) string {
	for name, version := range s.CompatToolVersions {
		for _, game := range version.Games {
			if id == game.ID {
				return name
			}
		}
	}

	return ""
}

func (s *Steam) IsValidVersion(version string) (bool, error) {
	for n, v := range s.CompatToolVersions {
		if version == n || version == v.ID {
			return true, nil
		}
	}

	fInfo, err := os.Stat(path.Join(s.root, "compatibilitytools.d", version))
	if err != nil {
		return false, err
	}

	if fInfo.IsDir() {
		return true, nil
	}

	return false, nil
}
