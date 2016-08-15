package bolt

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestList(t *testing.T) {
	db := newBoltDB(t)

	var err error
	key := []byte("userlist")
	bucket, _ := db.Bucket([]byte("1"))
	list := bucket.List(key)

	// insert a, b, c, d
	err = list.RPush([]byte("c"), []byte("d"))
	ensure.Nil(t, err)
	err = list.LPush([]byte("b"), []byte("a"))
	ensure.Nil(t, err)

	size, err := list.Len()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, size, int64(4))

	val, err := list.LPop()
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