package vdf2

import (
	"path"

	vdf "github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

type CompatToolMapping struct {
	Vdf
}

func (v *CompatToolMapping) Add(id, version string) {
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

func (v *CompatToolMapping) Update(id, version string) error {
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

func GetCompatToolMapping(steamRoot string) (*CompatToolMapping, error) {
	p := path.Join(steamRoot, "config", "config.vdf")

	n, err := parseTextConfig(p)
	if err != nil {
		return nil, err
	}

	key := []string{"InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping"}
	x, err := Lookup(n, key)
	if err != nil {
		return nil, err
	}

	return &CompatToolMapping{Vdf{n, x, p}}, nil
}
