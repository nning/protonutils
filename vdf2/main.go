package vdf2

import (
	"io/ioutil"
	"path"

	vdf "github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

type Vdf struct {
	Root *vdf.Node
	Node *vdf.Node
	Path string
}

func (v *Vdf) AddCompatToolMapping(id, version string) {
	addCompatToolMapping(v.Node, id, version)
}

func (v *Vdf) Save() error {
	out, err := v.Root.MarshalText()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(v.Path, out, 0600)
	if err != nil {
		return err
	}

	return nil
}

func Lookup(n *vdf.Node, x []string) (*vdf.Node, error) {
	y := n

	for _, key := range x {
		y = y.FirstByName(key)
		if y == nil {
			return nil, &steam.KeyNotFoundError{Name: key}
		}
	}

	return y, nil
}

func GetCompatToolMapping(steamRoot string) (*Vdf, error) {
	p := path.Join(steamRoot, "config", "config.vdf")

	n, err := parseTextConfig(p)
	if err != nil {
		return nil, err
	}

	key := []string{"Software", "Valve", "Steam", "CompatToolMapping"}
	x, err := Lookup(n, key)
	if err != nil {
		return nil, err
	}

	return &Vdf{n, x, p}, nil
}

func parseTextConfig(p string) (*vdf.Node, error) {
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var n vdf.Node
	err = n.UnmarshalText(in)

	return &n, nil
}

func addCompatToolMapping(n *vdf.Node, id, version string) {
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

	n.Append(&n0)
}
