package bolt

import (
	"bytes"
	"encoding/binary"
	"math"
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

// 使用二进制存储整形
func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
