package steam

import (
	"errors"
	"log/slog"
	"path/filepath"

	"github.com/BenLubar/vdf"
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
// ID. If version is empty, the mapping is removed, so that the game uses the
// default compatibility tool.
func (v *CompatToolMappingVdf) Update(id, version string) error {
	if version == "" {
		if id == "0" {
			return errors.New("cannot remove default compatibility tool mapping")
		}

		v.Node.FirstByName(id).Remove()
		return nil
	}

	x, err := Lookup(v.Node, []string{id, "name"})
	var keyNotFound *KeyNotFoundError

	if errors.As(err, &keyNotFound) {
		v.Node.FirstByName(id).Remove()
		v.Add(id, version)
	} else if err != nil {
		return err
	} else {
		x.SetString(version)
	}

	return nil
}

// GetDefaultID returns the ID of the compatibility tool that is configured as
// default (empty string if none)
func (v *CompatToolMappingVdf) GetDefaultID() string {
	return v.Node.FirstByName("0").FirstByName("name").String()
}

// ReadCompatTools reads compatibility tool mappings from VDF config and returns
// a CompatTools map with existing entries
func (v *CompatToolMappingVdf) ReadCompatTools() (CompatTools, error) {
	slog.Debug("CompatToolMappingVdf.ReadCompatTools")

	compatTools := make(CompatTools)

	for x := v.Node.FirstSubTree(); x != nil; x = x.NextSubTree() {
		id := x.Name()
		version := x.FirstByName("name").String()

		versionName, _ := v.Steam.GetCompatToolName(version)

		if id == "0" {
			if versionName == "" {
				versionName = "Default"
			} else {
				versionName += " (Default)"
			}

			compatTools.Add("", versionName)
			continue
		}

		game, isValid, err := v.Steam.GetGameData(id)
		if err != nil {
			return nil, err
		}

		if !isValid {
			continue
		}

		compatTools.Add(version, versionName)
		compatTools.AddGame(version, game)
	}

	return compatTools, nil
}

// IsValid checks whether a version is installed in the compatibility tools
// directory
func (v *CompatToolMappingVdf) IsValid(version string) bool {
	return v.Steam.CustomTools.Get(version) != nil
}

// GetCompatToolName returns human-readable name for compatibility tool
func (v *CompatToolMappingVdf) GetCompatToolName(version string) string {
	if version == "" {
		return v.GetDefaultID() + " (Default)"
	}

	return version
}

func (s *Steam) initCompatToolMapping() error {
	p := filepath.Join(s.Root, "config", "config.vdf")
	slog.Debug("steam.initCompatToolMapping", "path", p)

	n, err := ParseTextConfig(p)
	if err != nil {
		return err
	}

	key := []string{"Software", "Valve", "Steam", "CompatToolMapping"}
	x, err := Lookup(n, key)

	if err != nil {
		// Lookup is case-insensitive, so this only fails if the key path is
		// missing completely (e.g. on a fresh installation): create it
		x, err = ensurePath(n, key)
		if err != nil {
			return err
		}
	}

	s.CompatToolMapping = &CompatToolMappingVdf{Vdf{n, x, p, s}}

	return nil
}

// ensurePath returns the node at key path x below n and creates missing nodes
func ensurePath(n *vdf.Node, x []string) (*vdf.Node, error) {
	if n.FirstChild() == nil && n.String() != "" {
		return nil, errors.New("cannot create key path in value node")
	}

	y := n
	for _, key := range x {
		z := y.FirstByName(key)
		if z == nil {
			z = &vdf.Node{}
			z.SetName(key)
			y.Append(z)
		} else if z.FirstChild() == nil && z.String() != "" {
			return nil, &KeyNotFoundError{Name: key}
		}

		y = z
	}

	return y, nil
}
