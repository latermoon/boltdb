package boltdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

var ErrInvalidKeyFormat = errors.New("invalid key format includes +[],")

// Raw key:
// +key,type = value
// +name,s = "latermoon"

var (
	SEP = []byte{','}
	KEY = []byte{'+'} // Key Prefix
	SOK = []byte{'['} // Start of Key
	EOK = []byte{']'} // End of Key
)

// 字节范围
const (
	MINBYTE byte = 0
	MAXBYTE byte = math.MaxUint8
)

type ElemType byte

const (
	STRING    ElemType = 's'
	HASH      ElemType = 'h'
	LIST      ElemType = 'l'
	SORTEDSET ElemType = 'z'
	NONE      ElemType = '0'
)

func (e ElemType) String() string {
	switch byte(e) {
	case 's':
		return "string"
	case 'h':
		return "hash"
	case 'l':
		return "list"
	case 'z':
		return "sortedset"
	default:
		return "none"
	}
}

func rawKey(key []byte, t ElemType) []byte {
	return bytes.Join([][]byte{KEY, key, SEP, []byte{byte(t)}}, nil)
}

func verifyKey(key []byte) error {
	err := ErrInvalidKeyFormat
	if bytes.Contains(key, SEP) {
		return err
	} else if bytes.Contains(key, KEY) {
		return err
	} else if bytes.Contains(key, SOK) {
		return err
	} else if bytes.Contains(key, EOK) {
		return err
	}
	return nil
}

// itob returns an 8-byte big endian representation of v.
func itob(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func btoi(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
