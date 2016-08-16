package boltdb

import (
	"encoding/binary"
	"math"
)

// Score indicated that a number can be encoded to sorted []byte
// Score is use in SortedSet
// you can implement your own decode & encode function just like below
type Score []byte

// Int64ToScore ...
func Int64ToScore(i int64) Score {
	b := make([]byte, 9)
	// store sign in the first byte to keep the score order
	if i < 0 {
		b[0] = byte(0)
	} else {
		b[0] = byte(1)
	}
	binary.BigEndian.PutUint64(b[1:], uint64(i))
	return b
}

// ScoreToInt64 ...
func ScoreToInt64(b Score) int64 {
	return int64(binary.BigEndian.Uint64(b[1:]))
}

// Float64ToScore ...
func Float64ToScore(f float64) Score {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, float64ToUint64(f))
	return b
}

// ScoreToFloat64 ...
func ScoreToFloat64(b Score) float64 {
	return uint64ToFloat64(binary.BigEndian.Uint64(b))
}

// Copy from https://github.com/reborndb/qdb/blob/master/pkg/store/util.go
// We can not use lexicographically bytes comparison for negative and positive float directly.
// so here we will do a trick below.
func float64ToUint64(f float64) uint64 {
	u := math.Float64bits(f)
	if f >= 0 {
		u |= 0x8000000000000000
	} else {
		u = ^u
	}
	return u
}

func uint64ToFloat64(u uint64) float64 {
	if u&0x8000000000000000 > 0 {
		u &= ^uint64(0x8000000000000000)
	} else {
		u = ^u
	}
	return math.Float64frombits(u)
}
