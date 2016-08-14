package bolt

import (
	"github.com/boltdb/bolt"
)

type Bucket struct {
	db         *bolt.DB
	bucketName []byte
}

func (b *Bucket) Hash(key []byte) *Hash {
	return &Hash{bucket: b, key: key}
}

func (b *Bucket) List(key []byte) {

}

func (b *Bucket) SortedSet(key []byte) {

}

func (b *Bucket) Get(key []byte) ([]byte, error) {
	return b.rawGet(rawKey(key, STRING))
}

func (b *Bucket) Set(key, value []byte) error {
	return b.rawSet(rawKey(key, STRING), value)
}

func (b *Bucket) Drop(key []byte) {

}

func (b *Bucket) rawGet(key []byte) ([]byte, error) {
	var val []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(b.bucketName).Get(key)
		return nil
	})
	return val, err
}

func (b *Bucket) rawSet(key, value []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.bucketName).Put(key, value)
	})
}
