package vdf2

import (
	"path"

	"github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

type LocalConfigVdf struct {
	Vdf
}

func (v *LocalConfigVdf) GetViewedSteamPlay() ([]*steam.Game, error) {
	games := make([]*steam.Game, 0)
	var x *vdf.Node

	x = v.Node.FirstSubTree()

	for ; x != nil; x = x.NextChild() {
		id := x.Name()
		viewedSteamPlay := x.FirstByName("ViewedSteamPlay").String()

		if viewedSteamPlay != "1" {
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

func GetLocalConfig(s *steam.Steam) (*LocalConfigVdf, error) {
	p := path.Join(s.Root, "userdata", s.UID, "config", "localconfig.vdf")

	n, err := parseTextConfig(p)
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
