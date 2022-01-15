package vdf2

import (
	"os"
	"path"

	"github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

// CompatToolMappingVdf represents parsed VDF config for CompatToolMapping
type CompatToolMappingVdf struct {
	Vdf
}

// Add adds a new compatibility tool version mapping for a given app id
func (v *CompatToolMappingVdf) Add(appID, versionID string) {
	var n0 vdf.Node
	n0.SetName(appID)

	var n1 vdf.Node
	n1.SetName("name")
	n1.SetString(versionID)

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

// Update changes or adds a compatibility tool version mapping for a given app
// id
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

// ReadCompatTools reads compatibility tool mappings from VDF config and returns
// a CompatTools map with existing entries
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

		compatTools.Add(version, v.Steam.GetCompatToolName(version))
		compatTools.AddGame(version, game)
	}

	return compatTools, nil
}

// IsValid checks whether a version exists in the compatibility tools directory
func (v *CompatToolMappingVdf) IsValid(version string) bool {
	fInfo, err := os.Stat(path.Join(v.Steam.Root, "compatibilitytools.d", version))
	return err == nil && fInfo.IsDir()
}

// GetCompatToolMapping opens and parses config.vdf and returns a
// CompatToolMappingVdf containing the configurations
func GetCompatToolMapping(s *steam.Steam) (*CompatToolMappingVdf, error) {
	p := path.Join(s.Root, "config", "config.vdf")

	n, err := ParseTextConfig(p)
	if err != nil {
		return nil, err
	}

	key := []string{"Software", "Valve", "Steam", "CompatToolMapping"}
	x, err := Lookup(n, key)

	_, isKeyNotFoundError := err.(*steam.KeyNotFoundError)
	if err != nil && isKeyNotFoundError {
		key[2] = "steam"
		x, err = Lookup(n, key)
	}

	return &CompatToolMappingVdf{Vdf{n, x, p, s}}, err
}
