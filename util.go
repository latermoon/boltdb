package bolt

import (
	"bytes"
)

// Raw key:
// +key,type = value
// +name,s = "latermoon"

var (
	SEP = []byte{','}
	KEY = []byte{'+'} // Key Prefix
	SOK = []byte{'['} // Start of Key
	EOK = []byte{']'} // End of Key
)

type ElemType byte

const (
	STRING    ElemType = 's'
	HASH               = 'h'
	LIST               = 'l'
	SORTEDSET          = 'z'
	NONE               = '0'
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
