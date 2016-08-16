package bolt

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestList(t *testing.T) {
	db := newBoltDB(t)
	defer db.Close()

	var val []byte
	key := []byte("letter")
	bucket, _ := db.Bucket([]byte("1"))
	list, err := bucket.List(key)
	ensure.Nil(t, err)

	// insert a, b, c, d
	err = list.RPush([]byte("c"), []byte("d"))
	ensure.Nil(t, err)
	err = list.LPush([]byte("b"), []byte("a"))
	ensure.Nil(t, err)

	size, err := list.Len()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, size, int64(4))

	err = list.Range(0, 3, func(i int64, value []byte, quit *bool) {
		// log.Println(i, string(value))
	})
	ensure.Nil(t, err)

	val, err = list.Index(1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, val, []byte("b"))

	val, err = list.LPop()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, val, []byte("a"))

	size, err = list.Len()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, size, int64(3))

	val, err = list.RPop()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, val, []byte("d"))

	size, err = list.Len()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, size, int64(2))

	list.RPop()
	list.RPop()
	list.LPop()
	list.LPop()

	size, err = list.Len()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, size, int64(0))

	scan(db.db, []byte("1"), t)
}
