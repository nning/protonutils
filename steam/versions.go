package steam

import (
	"io"
	"os"
	"path"
	"sort"
)

// Version represents a compatibility tool version
type Version struct {
	ID        string
	Name      string
	Games     Games
	IsDefault bool
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

// ReadCompatToolVersions reads Proton versions and games from different Steam
// configs
func (s *Steam) ReadCompatToolVersions() error {
	x, err := s.getCompatToolMapping()
	if err != nil {
		return err
	}

	defID := x["0"].(mapLevel)["name"].(string)
	defName := s.getCompatToolName(defID) + " (Default)"

	for id, cfg := range x {
		vID := cfg.(mapLevel)["name"].(string)
		vName := s.getCompatToolName(vID)
		if vName == "" {
			vName = defName
		}

		_, err = s.addGame(vID, vName, id, false)
		if err != nil && err != io.EOF {
			return err
		}
	}

	x, err = s.getLocalConfig()
	_, isKeyNotFoundError := err.(*KeyNotFoundError)
	if err != nil && !isKeyNotFoundError {
		return err
	}

	for id, cfg := range x {
		v := cfg.(mapLevel)["ViewedSteamPlay"]
		if v == nil {
			continue
		}

		if !s.includesGameID(id) {
			_, err = s.addGame(defID, defName, id, true)
			if err != nil && err != io.EOF {
				return err
			}
		}
	}

	return nil
}

// GetGameVersion returns version ID (e.g. "proton_63") for a given game ID
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

// IsValidVersion returns whether a compatibility version is valid. version can
// either be a version ID (like "proton_63") or a human-readable name (like
// "Proton 6.3-8"). A version is valid if at least one game uses it or there is
// a install folder inside `compatibilitytools.d`.
func (s *Steam) IsValidVersion(version string) (bool, error) {
	if version == "" {
		return false, nil
	}

	for n, v := range s.CompatToolVersions {
		if version == n || version == v.ID {
			return true, nil
		}
	}

	fInfo, err := os.Stat(path.Join(s.Root, "compatibilitytools.d", version))
	if err != nil {
		return false, err
	}

	if fInfo.IsDir() {
		return true, nil
	}

	return false, nil
}
