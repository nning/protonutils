package vdf2

import (
	"bytes"
	"io/ioutil"
	"path"

	"github.com/BenLubar/vdf"
	"github.com/nning/protonutils/steam"
)

// AppInfoVdfApps represents apps found in appinfo VDF
type AppInfoVdfApps map[string]*vdf.Node

// AppInfoVdf wraps info for the appinfo VDF
type AppInfoVdf struct {
	Bytes []byte
	Apps  AppInfoVdfApps
	Path  string
	Steam *steam.Steam
}

// GetNextEntryStart returns the next offset to a appinfo binary VDF entry
// (starting at a given offset)
func (ai *AppInfoVdf) GetNextEntryStart(offset int) int {
	in := ai.Bytes
	needle := []byte("appinfo\x00")
	l := len(needle)

	for i := offset; i < len(in); i++ {
		if in[i] != needle[0] {
			continue
		}

		if bytes.Compare(in[i:i+l], needle) != 0 {
			continue
		}

		return i - 1
	}

	return -1
}

// ParseAppInfoEntry unmarshals `in` as binary VDF
func ParseAppInfoEntry(in []byte) (*vdf.Node, error) {
	var n vdf.Node
	err := n.UnmarshalBinary(in)
	return &n, err
}

// GetAppInfo loads appinfo VDF
func GetAppInfo(s *steam.Steam) (*AppInfoVdf, error) {
	p := path.Join(s.Root, "appcache", "appinfo.vdf")
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return &AppInfoVdf{
		Bytes: in,
		Apps:  make(AppInfoVdfApps),
		Path:  p,
		Steam: s,
	}, nil
}
