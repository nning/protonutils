package vdf2

import (
	"io/ioutil"

	"github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

// Vdf represents a parsed VDF file
type Vdf struct {
	Root  *vdf.Node
	Node  *vdf.Node
	Path  string
	Steam *steam.Steam
}

// Save saves a parsed VDF file back to disk
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
