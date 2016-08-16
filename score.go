package bolt

import "encoding/binary"

// Score indicated that a number can be encoded to sorted []byte
// Score is useful in SortedSet.Score
// you can implement your own decoder & encoder just like below
type Score []byte

// IntScore ...
func IntScore(i int64) Score {
	b := make([]byte, 9)
	// store sign in the first byte to keep the byte order
	if i < 0 {
		b[0] = byte(0)
	} else {
		b[0] = byte(1)
	}
	binary.BigEndian.PutUint64(b[1:], uint64(i))
	return b
}

// ScoreInt ...
func ScoreInt(b Score) int64 {
	return int64(binary.BigEndian.Uint64(b[1:]))
}
