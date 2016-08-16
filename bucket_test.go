package bolt

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestBucket(t *testing.T) {
	db := newBoltDB(t)

	b, err := db.Bucket([]byte("0"))
	ensure.Nil(t, err)

	hash, err := b.Hash([]byte("hash"))
	ensure.Nil(t, err)
	ensure.Nil(t, hash.Set([]byte("name"), []byte("latermoon")))

	// get a existing hash as a list
	_, err = b.List([]byte("hash"))
	ensure.DeepEqual(t, err, ErrWrongType)

	_, err = b.SortedSet([]byte("sortedset"))
	ensure.Nil(t, err)

}
