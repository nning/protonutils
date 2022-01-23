package steam2

import (
	"io/ioutil"

	"github.com/BenLubar/vdf"
)

// KeyNotFoundError is returned if key in VDF is not found
type KeyNotFoundError struct {
	Name string
}

func (e *KeyNotFoundError) Error() string {
	return "Key not found: " + e.Name
}

// Vdf represents a parsed VDF file
type Vdf struct {
	Root  *vdf.Node
	Node  *vdf.Node
	Path  string
	Steam *Steam
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
