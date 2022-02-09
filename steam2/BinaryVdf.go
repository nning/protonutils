package steam2

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/BenLubar/vdf"
)

// BinaryVdf wraps info for binary VDF
type BinaryVdf struct {
	Bytes []byte
	Path  string
}

func getAppIDNeedle(id string) ([]byte, error) {
	l := 10
	needle := make([]byte, 0, l)
	needle = append(needle, 'a', 'p', 'p', 'i', 'd', 0)

	n, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		n, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	slice := needle[6:l]
	binary.LittleEndian.PutUint32(slice, uint32(n))

	return needle[:l], nil
}

// GetNextEntryStart returns the next offset to a binary VDF entry beginning
// where needle was found (starting at given offset)
func (bVdf *BinaryVdf) GetNextEntryStart(offset, innerOffset int, needle []byte) int {
	if len(needle) == 0 {
		return -1
	}

	in := bVdf.Bytes
	l := len(needle)

	for i := offset; i+l < len(in); i++ {
		if in[i] != needle[0] {
			continue
		}

		if bytes.Compare(in[i:i+l], needle) != 0 {
			continue
		}

		return i - innerOffset
	}

	return -1
}

// GetNextEntryStartByID returns the next offset to a appinfo binary VDF entry
// by app id (starting at a given offset)
func (bVdf *BinaryVdf) GetNextEntryStartByID(offset, innerOffset int, id string) (int, error) {
	if id == "" || id == "0" {
		return -1, nil
	}

	needle, err := getAppIDNeedle(id)
	if err != nil {
		return -1, err
	}

	return bVdf.GetNextEntryStart(offset, innerOffset, needle), nil
}

// ParseBinaryVdf unmarshals `in` as binary VDF
func ParseBinaryVdf(in []byte) (*vdf.Node, error) {
	var n vdf.Node
	err := n.UnmarshalBinary(in)
	return &n, err
}
