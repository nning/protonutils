package steam

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"testing"

	vdf "github.com/wakeful-cloud/vdf"
)

func Test_ReadAppInfo(t *testing.T) {
	s, err := New(true)
	if err != nil {
		t.Error(err)
	}

	info, err := s.ReadAppInfo()
	if err != nil {
		t.Error(err)
	}
	if info.Magic != AppInfoMagic {
		t.Error("AppInfo Magic invalid")
	}
	if info.Universe != 1 {
		t.Error("Universe invalid")
	}
}

func Test_UseExternalLib(t *testing.T) {
	usr, _ := user.Current()
	file := path.Join(usr.HomeDir, ".steam", "root", "appcache", "appinfo.vdf")
	bytes, err := os.ReadFile(file)
	if err != nil {
		t.Error(err)
	}

	vdfMap, err := vdf.ReadVdf(bytes)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(vdfMap)
}
