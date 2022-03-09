package steam

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

// Lookup looks up a "key path" in a parsed VDF tree
func Lookup(n *vdf.Node, x []string) (*vdf.Node, error) {
	y := n

	for _, key := range x {
		y = y.FirstByName(key)
		if y == nil {
			return nil, &KeyNotFoundError{Name: key}
		}
	}

	return y, nil
}

// ParseTextConfig reads a file and parses it as text VDF
func ParseTextConfig(p string) (*vdf.Node, error) {
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var n vdf.Node
	err = n.UnmarshalText(in)

	return &n, nil
}
