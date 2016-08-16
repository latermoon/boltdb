package boltdb

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestHash(t *testing.T) {
	db := newBoltDB(t)
	defer db.Close()

	key := []byte("user:100422:profile")
	bucket, _ := db.Bucket([]byte("0"))
	hash, err := bucket.Hash(key)
	ensure.Nil(t, err)

	ensure.Nil(t, hash.Set([]byte("name"), []byte("latermoon")))
	ensure.Nil(t, hash.MSet([]byte("age"), []byte("28"), []byte("sex"), []byte("Male")))

	val, err := hash.Get([]byte("name"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, val, []byte("latermoon"))

	vals, err := hash.MGet([]byte("age"), []byte("sex"))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, vals[0], []byte("28"))
	ensure.DeepEqual(t, vals[1], []byte("Male"))

	keyVals, err := hash.GetAll()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, keyVals["age"], []byte("28"))
	ensure.DeepEqual(t, keyVals["name"], []byte("latermoon"))
	ensure.DeepEqual(t, keyVals["sex"], []byte("Male"))

	elemType, err := bucket.TypeOf(key)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, elemType, HASH)

	err = hash.Remove([]byte("name"), []byte("sex"))
	ensure.Nil(t, err)
	scan(db.db, []byte("0"), t)

	err = hash.Drop()
	ensure.Nil(t, err)
	scan(db.db, []byte("0"), t)
}
