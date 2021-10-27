package steam

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"os/user"
	"path"
	"strconv"
)

const AppInfoMagic = uint32(0x07_56_44_27)

type AppInfo struct {
	Magic    uint32
	Universe uint32
}

func getNeedle(appid string) ([]byte, error) {
	needle := make([]byte, 0, 10)
	needle = append(needle, 'a', 'p', 'p', 'i', 'd', 0)

	n, err := strconv.ParseInt(appid, 10, 32)
	if err != nil {
		return nil, err
	}

	slice := needle[6:10]
	binary.LittleEndian.PutUint32(slice, uint32(n))

	return needle[:10], nil
}

func (s *Steam) getAppInfoBuffer() (*AppInfo, io.Reader, error) {
	usr, _ := user.Current()
	file := path.Join(usr.HomeDir, ".steam", "root", "appcache", "appinfo.vdf")

	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	buf := bufio.NewReader(f)
	info := AppInfo{}
	err = binary.Read(buf, binary.LittleEndian, &info)
	if err != nil {
		return nil, nil, err
	}

	return &info, buf, nil
}
