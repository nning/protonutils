package steam

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"strconv"

	"github.com/BenLubar/vdf"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

// BinaryVdf wraps info for binary VDF
type BinaryVdf struct {
	Bytes []byte
	Path  string
}

// BinaryVdfTable is only used in debug mode to collect the count of duplicate
// parsing of binary VDF "snippets"
type BinaryVdfTable map[uint32]uint

var binaryVdfTable BinaryVdfTable

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

	log.Debug("BinaryVdf.GetNextEntryStart(", offset, innerOffset, needle, ")")

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
	if log.GetLevel() == log.DebugLevel {
		if binaryVdfTable == nil {
			binaryVdfTable = make(BinaryVdfTable)
		}

		crc := crc32.ChecksumIEEE(in)
		binaryVdfTable[crc]++
		log.Debug("steam.ParseBinaryVdf(", crc, "): ", humanize.Bytes(uint64(len(in))), ", ", binaryVdfTable[crc])
	}

	var n vdf.Node
	err := n.UnmarshalBinary(in)
	return &n, err
}
