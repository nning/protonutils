package steam

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

// Version represents a compatibility tool version
type Version struct {
	ID        string
	Name      string
	Games     Games
	IsDefault bool
	IsCustom  bool
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

// GetCompatToolName returns human-readable name of compatibility tool,
// for example: "proton_63" -> "Proton 6.3-8"
func (s *Steam) GetCompatToolName(shortName string) string {
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
	defName := s.GetCompatToolName(defID) + " (Default)"

	for id, cfg := range x {
		vID := cfg.(mapLevel)["name"].(string)
		vName := s.GetCompatToolName(vID)
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

	files, err := ioutil.ReadDir(path.Join(s.Root, "compatibilitytools.d"))
	if err != nil {
		return nil
	}

	for _, file := range files {
		id := file.Name()
		if !strings.HasPrefix(id, ".") {
			if s.CompatToolVersions[id] == nil {
				s.CompatToolVersions[id] = &Version{
					ID:       id,
					Name:     id,
					Games:    make(Games),
					IsCustom: true,
				}
			}
		}
	}

	return nil
}

// GetGameVersion returns Version struct for a given game ID
func (s *Steam) GetGameVersion(id string) *Version {
	for _, version := range s.CompatToolVersions {
		for _, game := range version.Games {
			if id == game.ID {
				return version
			}
		}
	}

	return nil
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

	return false, nil
}

// IsCustomVersion returns whether a version for a given versionID is manually
// installed in compatibilitytools.d
func (s *Steam) IsCustomVersion(versionID string) (bool, error) {
	if versionID == "" {
		return false, nil
	}

	fInfo, err := os.Stat(path.Join(s.Root, "compatibilitytools.d", versionID))
	if err != nil {
		return false, err
	}

	return fInfo.IsDir(), nil
}
