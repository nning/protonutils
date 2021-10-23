package steam

import (
	. "github.com/nning/list_proton_versions"
)

type CompatToolVersions map[string]Games

func (versions CompatToolVersions) IncludesGameId(id string) bool {
	for _, games := range versions {
		if games.IncludesId(id) {
			return true
		}
	}

	return false
}

func (s *Steam) IncludesGameId(id string) bool {
	return s.CompatToolVersions.IncludesGameId(id)
}

func (s *Steam) InitCompatToolVersions() {
	x, err := s.GetCompatToolMapping()
	PanicOnError(err)

	for id, cfg := range x {
		if id == "0" {
			continue
		}

		v := cfg.(MapLevel)["name"].(string)
		if v == "" {
			v = "Default"
		}

		s.AddGame(v, id)
	}

	x, err = s.GetLocalConfig()
	PanicOnError(err)

	for id, cfg := range x {
		v := cfg.(MapLevel)["ViewedSteamPlay"]
		if v == nil {
			continue
		}

		if !s.IncludesGameId(id) {
			s.AddGame("Default", id)
		}
	}
}
