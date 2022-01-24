package steam2

import (
	"bytes"
	"io/ioutil"
	"path"

	"github.com/BenLubar/vdf"
)

// AppInfoVdfApps represents apps found in appinfo VDF
type AppInfoVdfApps map[string]*vdf.Node

// AppInfoVdf wraps info for the appinfo VDF
type AppInfoVdf struct {
	Bytes []byte
	Apps  AppInfoVdfApps
	Path  string
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

func (ai *AppInfoVdf) GetName(id string) (string, error) {
	i := 0
	for {
		k := ai.GetNextEntryStart(i)
		if k < 0 {
			break
		}

		// TODO Do not parse every entry
		// TODO Return name and additional information like deck compatibility
		n, err := ParseAppInfoEntry(ai.Bytes[k:])
		if err != nil {
			return "", err

		}

		if id == n.FirstByName("appid").String() {
			return n.FirstByName("common").FirstByName("name").String(), nil
		}

		i = k + 2
	}

	return "", nil
}

// ParseAppInfoEntry unmarshals `in` as binary VDF
func ParseAppInfoEntry(in []byte) (*vdf.Node, error) {
	var n vdf.Node
	err := n.UnmarshalBinary(in)
	return &n, err
}

// GetAppInfo loads appinfo VDF
func (s *Steam) InitAppInfo() error {
	p := path.Join(s.Root, "appcache", "appinfo.vdf")
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	s.AppInfo = &AppInfoVdf{
		Bytes: in,
		Apps:  make(AppInfoVdfApps),
		Path:  p,
	}

	return nil
}
