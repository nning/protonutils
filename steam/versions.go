package steam

import (
	"sort"
)

// CompatToolVersions maps Proton versions to games
type CompatToolVersions map[string]Games

func (versions CompatToolVersions) includesGameID(id string) bool {
	for _, games := range versions {
		if games.includesID(id) {
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

// ReadCompatToolVersions reads Proton versions and games from different Steam configs
func (s *Steam) ReadCompatToolVersions() error {
	x, err := s.getCompatToolMapping()
	if err != nil {
		return err
	}

	def := x["0"].(mapLevel)["name"].(string) + " (Default)"

	for id, cfg := range x {
		v := cfg.(mapLevel)["name"].(string)
		if v == "" {
			v = def
		}

		_, err = s.addGame(v, id)
		if err != nil {
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
			_, err = s.addGame(def, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
