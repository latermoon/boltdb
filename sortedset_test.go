package bolt

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestIntScore(t *testing.T) {
	a, b, c := IntScore(-1), IntScore(0), IntScore(1)
	ensure.DeepEqual(t, ScoreInt(a), int64(-1))
	ensure.DeepEqual(t, ScoreInt(b), int64(0))
	ensure.DeepEqual(t, ScoreInt(c), int64(1))
}

func TestSortedSet(t *testing.T) {
	db := newBoltDB(t)
	defer db.Close()

	// var val []byte
	key := []byte("users")
	bucket, _ := db.Bucket([]byte("2"))
	zset, err := bucket.SortedSet(key)
	ensure.Nil(t, err)

	added, err := zset.Add(IntScore(-1), []byte("a"), IntScore(0), []byte("b"), IntScore(1), []byte("c"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 3)

	score, err := zset.Score([]byte("b"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, ScoreInt(score), int64(0))

	added, err = zset.Add(IntScore(100), []byte("d"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 1)

	added, err = zset.Add(IntScore(200), []byte("b"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 0) // no new record(s)

	err = zset.RangeByScore(IntScore(-1), IntScore(200), func(i int64, score Score, member []byte, quit *bool) {
		t.Log("Range", i, ScoreInt(score), string(member))
	})
	ensure.Nil(t, err)

	err = zset.RevRangeByScore(IntScore(100), IntScore(-1), func(i int64, score Score, member []byte, quit *bool) {
		t.Log("RevRange", i, ScoreInt(score), string(member))
	})
	ensure.Nil(t, err)

	removed, err := zset.Remove([]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, removed, 4)

	scan(db.db, []byte("2"), t)
}
