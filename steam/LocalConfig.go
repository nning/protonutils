package steam

import (
	"log/slog"
	"path/filepath"
)

// LocalConfigVdf represents parsed VDF config for app config from
// localconfig.vdf
type LocalConfigVdf struct {
	Vdf
}

// GetGames returns a slice of games read from localconfig.vdf (Steam tools like
// Proton or the Steam Linux Runtime are skipped).
// The `enableViewedSteamPlay` argument controls whether the slice only contains
// games for which the user confirmed the Steam Play disclaimer
func (v *LocalConfigVdf) GetGames(enableViewedSteamPlay bool) ([]*Game, error) {
	slog.Debug("LocalConfigVdf.GetGames")

	games := make([]*Game, 0)

	for x := v.Node.FirstSubTree(); x != nil; x = x.NextSubTree() {
		id := x.Name()

		if enableViewedSteamPlay && x.FirstByName("ViewedSteamPlay").String() != "1" {
			continue
		}

		if v.Steam.AppInfo.IsTool(id) {
			continue
		}

		game, isValid, err := v.Steam.GetGameData(id)
		if err != nil {
			return nil, err
		}

		if !isValid {
			continue
		}

		games = append(games, game)
	}

	return games, nil
}

func (s *Steam) initLocalConfig() error {
	p := filepath.Join(s.Root, "userdata", s.UID, "config", "localconfig.vdf")
	slog.Debug("steam.initLocalConfig", "path", p)

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	key := []string{"Software", "Valve", "Steam", "apps"}
	x, err := Lookup(n, key)
	if err != nil {
		return err
	}

	s.LocalConfig = &LocalConfigVdf{Vdf{n, x, p, s}}

	return nil
}
