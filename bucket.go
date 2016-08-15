package bolt

import (
	"bytes"

	"github.com/boltdb/bolt"
	// "errors"
)

type Bucket struct {
	db         *bolt.DB
	bucketName []byte
}

func (b *Bucket) Hash(key []byte) *Hash {
	return &Hash{bucket: b, key: key}
}

func (b *Bucket) List(key []byte) *List {
	return &List{bucket: b, key: key}
}

func (b *Bucket) SortedSet(key []byte) *SortedSet {
	return &SortedSet{bucket: b, key: key}
}

func (b *Bucket) TypeOf(key []byte) (ElemType, error) {
	elemType := NONE
	err := b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(b.bucketName).Cursor()
		prefix := bytes.Join([][]byte{KEY, key, SEP}, nil)
		if k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix) {
			t := bytes.TrimPrefix(k, prefix)
			elemType = ElemType(t[0])
		}
		return nil
	})
	return elemType, err
}

func (b *Bucket) Get(key []byte) ([]byte, error) {
	var val []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(b.bucketName).Get(rawKey(key, STRING))
		return nil
	})
	return val, err
}

func (b *Bucket) Set(key, value []byte) error {
	return b.db.Batch(func(tx *bolt.Tx) error {
		return tx.Bucket(b.bucketName).Put(rawKey(key, STRING), value)
	})
}
