package vdf2

import (
	"io/ioutil"

	vdf "github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

type Vdf struct {
	Root *vdf.Node
	Node *vdf.Node
	Path string
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

	if y.Name() != x[0] {
		return nil, &steam.KeyNotFoundError{Name: x[0]}
	}

	for _, key := range x[1:] {
		y = y.FirstByName(key)
		if y == nil {
			return nil, &steam.KeyNotFoundError{Name: key}
		}
	}

	return y, nil
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
