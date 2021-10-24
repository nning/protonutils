package steam

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

func (s *Steam) InitCompatToolVersions(user string) error {
	x, err := s.GetCompatToolMapping()
	if err != nil {
		return err
	}

	def := x["0"].(MapLevel)["name"].(string) + " (Default)"

	for id, cfg := range x {
		v := cfg.(MapLevel)["name"].(string)
		if v == "" {
			v = def
		}

		_, err = s.AddGame(v, id)
		if err != nil {
			return err
		}
	}

	if user != "" {
		user, err = s.UserToId32(user)
		if err != nil {
			return err
		}
	}

	x, err = s.GetLocalConfig(user)
	if err != nil {
		return err
	}

	for id, cfg := range x {
		v := cfg.(MapLevel)["ViewedSteamPlay"]
		if v == nil {
			continue
		}

		if !s.IncludesGameId(id) {
			_, err = s.AddGame(def, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
