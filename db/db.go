package db

import (
	"errors"
	"github.com/boltdb/bolt"
)

var ErrBadRequest = errors.New("Bad request")

type BoltService struct {
	db     *bolt.DB
	bucket []byte
}

func New(dbpath string) (*BoltService, error) {
	db, err := bolt.Open(dbpath, 0644, nil)
	if err != nil {
		return nil, err
	}

	b := &BoltService{
		db:     db,
		bucket: []byte("0"), // default bucket
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b.bucket)
		return err
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BoltService) Close() error {
	return b.db.Close()
}

func (b *BoltService) SelectBucket(name []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
	if err == nil {
		b.bucket = name
	}
	return err
}

func (b *BoltService) Set(key, value []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.bucket).Put(key, value)
	})
	return err
}

func (b *BoltService) MSet(kvs ...[]byte) error {
	if len(kvs)%2 != 0 {
		return ErrBadRequest
	}
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		for i := 0; i < len(kvs); i += 2 {
			key, value := kvs[i], kvs[i+1]
			if err := bucket.Put(key, value); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (b *BoltService) Get(key []byte) ([]byte, error) {
	var val []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(b.bucket).Get(key)
		return nil
	})
	return val, err
}

// TODO use MultiBulks better
func (b *BoltService) MGet(keys ...[]byte) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := b.db.View(func(tx *bolt.Tx) error {
		for _, key := range keys {
			val := tx.Bucket(b.bucket).Get(key)
			if val != nil {
				result[string(key)] = string(val)
			} else {
				result[string(key)] = nil
			}
		}
		return nil
	})
	return result, err
}

func (b *BoltService) Delete(keys ...[]byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		for _, key := range keys {
			if err := tx.Bucket(b.bucket).Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}
