package steam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"math"

	"github.com/BenLubar/vdf"
)

// Type bytes of binary VDF (KeyValues) nodes
const (
	kvTypeNone    = 0x00 // sub tree
	kvTypeString  = 0x01
	kvTypeInt32   = 0x02
	kvTypeFloat32 = 0x03
	kvTypePointer = 0x04
	kvTypeWString = 0x05
	kvTypeColor   = 0x06
	kvTypeUint64  = 0x07
	kvTypeEnd     = 0x08
	kvTypeInt64   = 0x0a
	kvTypeEndAlt  = 0x0b
)

// maxBinaryVdfDepth limits nesting to protect against corrupt input
const maxBinaryVdfDepth = 256

// BinaryVdfError is returned if binary VDF data cannot be parsed
type BinaryVdfError struct {
	Offset int
	Msg    string
}

func (e *BinaryVdfError) Error() string {
	return fmt.Sprintf("binary vdf: %v at offset %d", e.Msg, e.Offset)
}

type binaryVdfReader struct {
	data []byte
	pos  int

	// If keys is not nil, node names are stored as uint32 indices into this
	// table instead of null-terminated strings (appinfo.vdf v29)
	keys []string
}

// ParseBinaryVdf unmarshals `in` as binary VDF and returns the first top-level
// node (further top-level nodes are reachable as its siblings)
func ParseBinaryVdf(in []byte) (*vdf.Node, error) {
	return parseBinaryVdf(in, nil)
}

func parseBinaryVdf(in []byte, keys []string) (*vdf.Node, error) {
	r := &binaryVdfReader{data: in, keys: keys}

	var root vdf.Node
	err := r.readChildren(&root, 0)
	if err != nil {
		return nil, err
	}

	if root.FirstChild() == nil {
		return nil, &BinaryVdfError{0, "no nodes found"}
	}

	return root.FirstChild(), nil
}

func (r *binaryVdfReader) err(msg string) error {
	return &BinaryVdfError{r.pos, msg}
}

func (r *binaryVdfReader) need(n int) error {
	if n < 0 || r.pos+n > len(r.data) {
		return r.err("unexpected end of data")
	}

	return nil
}

func (r *binaryVdfReader) readByte() (byte, error) {
	if err := r.need(1); err != nil {
		return 0, err
	}

	b := r.data[r.pos]
	r.pos++
	return b, nil
}

func (r *binaryVdfReader) readUint16() (uint16, error) {
	if err := r.need(2); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return v, nil
}

func (r *binaryVdfReader) readUint32() (uint32, error) {
	if err := r.need(4); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return v, nil
}

func (r *binaryVdfReader) readUint64() (uint64, error) {
	if err := r.need(8); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint64(r.data[r.pos:])
	r.pos += 8
	return v, nil
}

func (r *binaryVdfReader) readCString() (string, error) {
	i := bytes.IndexByte(r.data[r.pos:], 0)
	if i < 0 {
		return "", r.err("unterminated string")
	}

	s := string(r.data[r.pos : r.pos+i])
	r.pos += i + 1
	return s, nil
}

func (r *binaryVdfReader) readKey() (string, error) {
	if r.keys == nil {
		return r.readCString()
	}

	i, err := r.readUint32()
	if err != nil {
		return "", err
	}

	if int(i) >= len(r.keys) {
		return "", r.err(fmt.Sprintf("key index %d out of range", i))
	}

	return r.keys[i], nil
}

// readChildren reads nodes and appends them to parent until an end marker (or,
// on top-level, the end of data) is reached
func (r *binaryVdfReader) readChildren(parent *vdf.Node, depth int) error {
	if depth > maxBinaryVdfDepth {
		return r.err("nesting too deep")
	}

	for {
		if depth == 0 && r.pos >= len(r.data) {
			return nil
		}

		t, err := r.readByte()
		if err != nil {
			return err
		}

		if t == kvTypeEnd || t == kvTypeEndAlt {
			return nil
		}

		name, err := r.readKey()
		if err != nil {
			return err
		}

		n := &vdf.Node{}
		n.SetName(name)

		switch t {
		case kvTypeNone:
			parent.Append(n)
			err = r.readChildren(n, depth+1)
			if err != nil {
				return err
			}
			continue
		case kvTypeString:
			var v string
			v, err = r.readCString()
			n.SetString(v)
		case kvTypeInt32:
			var v uint32
			v, err = r.readUint32()
			n.SetInt(int32(v))
		case kvTypeFloat32:
			var v uint32
			v, err = r.readUint32()
			n.SetFloat(math.Float32frombits(v))
		case kvTypePointer:
			var v uint32
			v, err = r.readUint32()
			n.SetPtr(v)
		case kvTypeWString:
			var l uint16
			l, err = r.readUint16()
			if err != nil {
				return err
			}
			v := make([]uint16, l)
			for i := range v {
				v[i], err = r.readUint16()
				if err != nil {
					return err
				}
			}
			n.SetWString(v)
		case kvTypeColor:
			err = r.need(4)
			if err == nil {
				b := r.data[r.pos : r.pos+4]
				n.SetColor(color.NRGBA{R: b[0], G: b[1], B: b[2], A: b[3]})
				r.pos += 4
			}
		case kvTypeUint64, kvTypeInt64:
			var v uint64
			v, err = r.readUint64()
			n.SetUint64(v)
		default:
			r.pos--
			return r.err(fmt.Sprintf("unknown type 0x%02x", t))
		}

		if err != nil {
			return err
		}

		parent.Append(n)
	}
}
