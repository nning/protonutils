package steam

import (
	"bufio"
	"encoding/binary"
	"os"
	"os/user"
	"path"
)

const AppInfoMagic = uint32(0x07_56_44_27)

type AppInfo struct {
	Magic    uint32
	Universe uint32
	// Apps     []AppSection
}

// type AppSection struct {
// 	AppId        uint32
// 	Size         uint32
// 	InfoState    uint32
// 	LastUpdated  uint32
// 	PicsToken    uint64
// 	Sha1         [20]byte
// 	ChangeNumber uint32
// }

// See https://github.com/SteamDatabase/SteamAppInfo/blob/master/SteamAppInfoParser/AppInfo.cs
func (s *Steam) ReadAppInfo() (*AppInfo, error) {
	usr, _ := user.Current()
	file := path.Join(usr.HomeDir, ".steam", "root", "appcache", "appinfo.vdf")

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(f)
	info := AppInfo{}
	err = binary.Read(buf, binary.LittleEndian, &info)
	if err != nil {
		return nil, err
	}

	// b, err := buf.ReadByte()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(b) // is 05, read position is OK

	return &info, nil
}
