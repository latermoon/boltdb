package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/facebookgo/ensure"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	db := newBoltDB(t)

	var err error
	bucket, _ := db.Bucket([]byte("0"))
	hash := bucket.Hash([]byte("user:100422:profile"))

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
	ensure.DeepEqual(t, keyVals[0], []byte("age"))
	ensure.DeepEqual(t, keyVals[1], []byte("28"))
	ensure.DeepEqual(t, keyVals[2], []byte("name"))
	ensure.DeepEqual(t, keyVals[3], []byte("latermoon"))
	ensure.DeepEqual(t, keyVals[4], []byte("sex"))
	ensure.DeepEqual(t, keyVals[5], []byte("Male"))

	err = hash.Remove([]byte("name"), []byte("sex"), []byte("age"))
	ensure.Nil(t, err)

	scan(db.db, []byte("0"), t)
}

func scan(db *bolt.DB, bucket []byte, t *testing.T) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Printf("%s  %s\n", k, v)
		}
		return nil
	})
}
