package bolt

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestSortedSet(t *testing.T) {
	db := newBoltDB(t)

	var err error
	// var val []byte
	key := []byte("users")
	bucket, _ := db.Bucket([]byte("2"))
	zset := bucket.SortedSet(key)

	added, err := zset.Add([]byte("1"), []byte("a"), []byte("2"), []byte("b"), []byte("3"), []byte("c"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 3)

	score, err := zset.Score([]byte("b"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, score, []byte("2"))

	added, err = zset.Add([]byte("4"), []byte("d"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 1)

	added, err = zset.Add([]byte("200"), []byte("b"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, added, 0) // no new record(s)

	err = zset.RangeByScore([]byte("1"), []byte("z"), func(i int64, score, member []byte, quit *bool) {
		// log.Println(i, string(score), string(member))
	})
	ensure.Nil(t, err)

	err = zset.RevRangeByScore([]byte("3"), []byte("1"), func(i int64, score, member []byte, quit *bool) {
		// log.Println(i, string(score), string(member))
	})
	ensure.Nil(t, err)

	removed, err := zset.Remove([]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, removed, 4)

	scan(db.db, []byte("2"), t)
}
