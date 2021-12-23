package vdf2

import (
	"path"

	vdf "github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

type CompatTools map[string]*CompatTool

type CompatTool = steam.Version

type CompatToolMappingVdf struct {
	Vdf
}

type CompatToolsVdf []*vdf.Node

func (c CompatTools) IsValid(v string) bool {
	for id := range c {
		if id == v {
			return true
		}
	}

	return false
}

func (v *CompatToolMappingVdf) Add(id, version string) {
	var n0 vdf.Node
	n0.SetName(id)

	var n1 vdf.Node
	n1.SetName("name")
	n1.SetString(version)

	var n2 vdf.Node
	n2.SetName("config")
	n2.SetString("")

	var n3 vdf.Node
	n3.SetName("Priority")
	n3.SetString("250")

	n0.Append(&n1)
	n0.Append(&n2)
	n0.Append(&n3)

	v.Node.Append(&n0)
}

func (v *CompatToolMappingVdf) Update(id, version string) error {
	x, err := Lookup(v.Node, []string{id, "name"})
	_, isKeyNotFoundError := err.(*steam.KeyNotFoundError)

	if isKeyNotFoundError {
		v.Add(id, version)
	} else if err != nil {
		return err
	} else {
		x.SetString(version)
	}

	return nil
}

func (v *CompatToolMappingVdf) ReadCompatTools() (CompatTools, error) {
	compatTools := make(CompatTools)
	var x *vdf.Node

	x = v.Node.FirstSubTree()

	for ; x != nil; x = x.NextChild() {
		id := x.Name()
		version := x.FirstByName("name").String()

		game, isValid, err := v.Steam.GetGameData(id)
		if err != nil {
			return nil, err
		}

		if !isValid {
			continue
		}

		if compatTools[version] == nil {
			compatTools[version] = &CompatTool{
				ID:        version,
				Name:      v.Steam.GetCompatToolName(version),
				IsDefault: false,
				Games:     make(steam.Games),
			}
		}

		compatTools[version].Games[game.Name] = game
	}

	return compatTools, nil
}

func GetCompatToolMapping(s *steam.Steam) (*CompatToolMappingVdf, error) {
	p := path.Join(s.Root, "config", "config.vdf")

	n, err := parseTextConfig(p)
	if err != nil {
		return nil, err
	}

	key := []string{"InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping"}
	x, err := Lookup(n, key)
	if err != nil {
		return nil, err
	}

	return &CompatToolMappingVdf{Vdf{n, x, p, s}}, nil
}
