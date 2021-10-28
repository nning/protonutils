package steam

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"os/user"
	"path"
	"strconv"
)

const appInfoMagic = uint32(0x07_56_44_27)

type appInfo struct {
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

func (s *Steam) getAppInfoBuffer() (*appInfo, *bufio.Reader, error) {
	usr, _ := user.Current()
	file := path.Join(usr.HomeDir, ".steam", "root", "appcache", "appinfo.vdf")

	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	buf := bufio.NewReader(f)
	info := appInfo{}
	err = binary.Read(buf, binary.LittleEndian, &info)
	if err != nil {
		return nil, nil, err
	}

	return &info, buf, nil
}

func findNeedleInBuffer(buf *bufio.Reader, needle []byte) (string, error) {
	l := len(needle)

	for {
		b, err := buf.ReadByte()
		if err != nil {
			return "", err
		}

		if b != needle[0] {
			continue
		}

		hay, err := buf.Peek(l - 1)
		if err != nil {
			return "", err
		}

		if bytes.Compare(hay, needle[1:]) != 0 {
			continue
		}

		_, err = buf.Discard(l - 1)
		if err != nil {
			return "", err
		}

		newNeedle := []byte{'n', 'a', 'm', 'e', 0}
		if bytes.Compare(needle, newNeedle) != 0 {
			return findNeedleInBuffer(buf, newNeedle)
		}

		s, err := buf.ReadBytes(0)
		if err != nil {
			return "", err
		}

		return string(s[:len(s)-1]), nil
	}
}

func (s *Steam) findNameInAppInfo(id string) (string, error) {
	_, buf, err := s.getAppInfoBuffer()
	if err != nil {
		return "", err
	}

	needle, err := getNeedle(id)
	if err != nil {
		return "", err
	}

	return findNeedleInBuffer(buf, needle)
}
