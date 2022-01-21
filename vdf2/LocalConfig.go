package vdf2

import (
	"path"

	"github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

// LocalConfigVdf represents parsed VDF config for app config from
// localconfig.vdf
type LocalConfigVdf struct {
	Vdf
}

// GetViewedSteamPlay returns a slice of games for which the user confirmed the
// Steam Play disclaimer
func (v *LocalConfigVdf) GetViewedSteamPlay() ([]*Game, error) {
	games := make([]*Game, 0)
	var x *vdf.Node

	x = v.Node.FirstSubTree()

	for ; x != nil; x = x.NextChild() {
		id := x.Name()
		viewedSteamPlay := x.FirstByName("ViewedSteamPlay").String()

		if viewedSteamPlay != "1" {
			continue
		}

		game, isValid, err := GetGameData(v.Steam, id)
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

// GetLocalConfig reads and parses localconfig.vdf and returns a LocalConfigVdf
func GetLocalConfig(s *steam.Steam) (*LocalConfigVdf, error) {
	p := path.Join(s.Root, "userdata", s.UID, "config", "localconfig.vdf")

	n, err := ParseTextConfig(p)
	if err != nil {
		return nil, err
	}

	key := []string{"Software", "Valve", "Steam", "apps"}
	x, err := Lookup(n, key)
	if err != nil {
		return nil, err
	}

	return &LocalConfigVdf{Vdf{n, x, p, s}}, nil
}
