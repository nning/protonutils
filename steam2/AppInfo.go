package steam2

import (
	"io/ioutil"
	"path"
)

// InnerOffsetAppInfo sets byte count before appid match in appinfo
const InnerOffsetAppInfo = 10

func (s *Steam) initAppInfo() error {
	p := path.Join(s.Root, "appcache", "appinfo.vdf")
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	s.AppInfo = &BinaryVdf{
		Bytes: in,
		Path:  p,
	}

	return nil
}
