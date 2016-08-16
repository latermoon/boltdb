package bolt

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

var ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

type Bucket struct {
	db         *bolt.DB
	bucketName []byte
}

func (b *Bucket) Hash(key []byte) (*Hash, error) {
	if err := b.ensureType(key, HASH); err != nil {
		return nil, err
	}
	return &Hash{bucket: b, key: key}, nil
}

func (b *Bucket) List(key []byte) (*List, error) {
	if err := b.ensureType(key, LIST); err != nil {
		return nil, err
	}
	return &List{bucket: b, key: key}, nil
}

func (b *Bucket) SortedSet(key []byte) (*SortedSet, error) {
	if err := b.ensureType(key, SORTEDSET); err != nil {
		return nil, err
	}
	return &SortedSet{bucket: b, key: key}, nil
}

func (b *Bucket) ensureType(key []byte, elemType ElemType) error {
	if t, err := b.TypeOf(key); err != nil {
		return err
	} else if t != NONE && t != elemType {
		return ErrWrongType
	}
	return nil
}

func (b *Bucket) TypeOf(key []byte) (ElemType, error) {
	elemType := NONE
	err := b.View(func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()
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
	err := b.View(func(bucket *bolt.Bucket) error {
		val = bucket.Get(rawKey(key, STRING))
		return nil
	})
	return val, err
}

func (b *Bucket) Set(key, value []byte) error {
	return b.Batch(func(bucket *bolt.Bucket) error {
		return bucket.Put(rawKey(key, STRING), value)
	})
}

// View make bolt.DB.View(Tx){tx.Bucket(...)} to View(bolt.Bucket)
func (b *Bucket) View(fn func(*bolt.Bucket) error) error {
	return b.db.View(func(tx *bolt.Tx) error {
		return fn(tx.Bucket(b.bucketName))
	})
}

// Batch ...
func (b *Bucket) Batch(fn func(*bolt.Bucket) error) error {
	return b.db.Batch(func(tx *bolt.Tx) error {
		return fn(tx.Bucket(b.bucketName))
	})
}
