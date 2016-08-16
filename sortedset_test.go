package boltdb

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestScore(t *testing.T) {
	a, b, c := Int64ToScore(-1), Int64ToScore(0), Int64ToScore(1)
	ensure.DeepEqual(t, ScoreToInt64(a), int64(-1))
	ensure.DeepEqual(t, ScoreToInt64(b), int64(0))
	ensure.DeepEqual(t, ScoreToInt64(c), int64(1))

	x, y, z := Float64ToScore(-1.5), Float64ToScore(0), Float64ToScore(1.5)
	ensure.DeepEqual(t, ScoreToFloat64(x), float64(-1.5))
	ensure.DeepEqual(t, ScoreToFloat64(y), float64(0))
	ensure.DeepEqual(t, ScoreToFloat64(z), float64(1.5))
}

func TestSortedSet(t *testing.T) {
	db := newBoltDB(t)
	defer db.Close()

	// var val []byte
	key := []byte("users")
	bucket, _ := db.Bucket([]byte("2"))
	zset, err := bucket.SortedSet(key)
	ensure.Nil(t, err)

	added, err := zset.Add(Int64ToScore(-1), []byte("a"), Int64ToScore(0), []byte("b"), Int64ToScore(1), []byte("c"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 3)

	score, err := zset.Score([]byte("b"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, ScoreToInt64(score), int64(0))

	added, err = zset.Add(Int64ToScore(100), []byte("d"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 1)

	added, err = zset.Add(Int64ToScore(200), []byte("b"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 0) // no new record(s)

	err = zset.RangeByScore(Int64ToScore(-1), Int64ToScore(200), func(i int64, score Score, member []byte, quit *bool) {
		t.Log("Range", i, ScoreToInt64(score), string(member))
	})
	ensure.Nil(t, err)

	err = zset.RevRangeByScore(Int64ToScore(100), Int64ToScore(-1), func(i int64, score Score, member []byte, quit *bool) {
		t.Log("RevRange", i, ScoreToInt64(score), string(member))
	})
	ensure.Nil(t, err)

	removed, err := zset.Remove([]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, removed, 4)

	scan(db.db, []byte("2"), t)
}
