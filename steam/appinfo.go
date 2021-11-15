package steam

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

func getAppIDNeedle(appid string) ([]byte, error) {
	l := 10
	needle := make([]byte, 0, l)
	needle = append(needle, 'a', 'p', 'p', 'i', 'd', 0)

	n, err := strconv.ParseInt(appid, 10, 32)
	if err != nil {
		n, err = strconv.ParseInt(appid, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	slice := needle[6:l]
	binary.LittleEndian.PutUint32(slice, uint32(n))

	return needle[:l], nil
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

func (s *Steam) getShortcutsBuffer() (*bufio.Reader, error) {
	usr, _ := user.Current()
	file := path.Join(usr.HomeDir, ".steam", "root", "userdata", s.uid, "config", "shortcuts.vdf")

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(f)

	return buf, nil
}

// findNeedleInBuffer searches in buf:
//   * First search for needle1
//   * If needle1 has been found, search for needle2
//   * If needle2 has been found, return bytes right after (until nullbyte)
// Will check lookAhead bytes, use -1 for no limit.
func findNeedleInBuffer(buf *bufio.Reader, needle1, needle2 []byte, lookAhead int) (string, error) {
	l := len(needle1)

	for {
		b, err := buf.ReadByte()
		if err != nil {
			return "", err
		}

		lookAhead--
		if lookAhead == 0 {
			return "", io.EOF
		}

		if b != needle1[0] {
			continue
		}

		hay, err := buf.Peek(l - 1)
		if err != nil {
			return "", err
		}

		if bytes.Compare(hay, needle1[1:]) != 0 {
			continue
		}

		_, err = buf.Discard(l - 1)
		if err != nil {
			return "", err
		}

		if len(needle2) > 0 {
			return findNeedleInBuffer(buf, needle2, nil, 1024)
		}

		s, err := buf.ReadBytes(0)
		if err != nil {
			return "", err
		}

		return string(s[:len(s)-1]), nil
	}
}

func (s *Steam) findNameInAppInfo(id string) (string, error) {
	if id == "0" {
		return "", nil
	}

	_, buf, err := s.getAppInfoBuffer()
	if err != nil {
		return "", err
	}

	needle1, err := getAppIDNeedle(id)
	if err != nil {
		return "", err
	}

	fmt.Println(id, needle1)

	needle2 := []byte{'n', 'a', 'm', 'e', 0}

	return findNeedleInBuffer(buf, needle1, needle2, -1)
}

func (s *Steam) findNameInShortcuts(id string) (string, error) {
	buf, err := s.getShortcutsBuffer()
	if err != nil {
		return "", err
	}

	needle1, err := getAppIDNeedle(id)
	if err != nil {
		return "", err
	}

	needle2 := []byte{'A', 'p', 'p', 'N', 'a', 'm', 'e', 0}

	return findNeedleInBuffer(buf, needle1, needle2, -1)
}
