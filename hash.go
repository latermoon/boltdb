package bolt

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

// Hash ...
// 	+key,h = ""
// 	h[key]name = "latermoon"
// 	h[key]age = "27"
// 	h[key]sex = "Male"
type Hash struct {
	bucket *Bucket
	key    []byte
}

func (h *Hash) Get(field []byte) ([]byte, error) {
	var val []byte
	err := h.bucket.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(h.bucket.bucketName)
		val = b.Get(h.fieldKey(field))
		return nil
	})
	return val, err
}

func (h *Hash) MGet(fields ...[]byte) ([][]byte, error) {
	vals := make([][]byte, 0, len(fields))
	err := h.bucket.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(h.bucket.bucketName)
		for _, field := range fields {
			val := b.Get(h.fieldKey(field))
			vals = append(vals, val)
		}
		return nil
	})
	return vals, err
}

// GetAll ...
// TODO: add seek & max return count
func (h *Hash) GetAll() ([][]byte, error) {
	keyVals := make([][]byte, 0)
	err := h.bucket.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(h.bucket.bucketName).Cursor()
		prefix := h.fieldPrefix()
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			keyVals = append(keyVals, h.fieldInKey(k), v)
		}
		return nil
	})
	return keyVals, err
}

func (h *Hash) Set(field, value []byte) error {
	return h.MSet(field, value)
}

func (h *Hash) MSet(fieldVals ...[]byte) error {
	if len(fieldVals) == 0 || len(fieldVals)%2 != 0 {
		return errors.New("invalid field value pairs")
	}

	return h.bucket.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(h.bucket.bucketName)
		for i := 0; i < len(fieldVals); i += 2 {
			field, val := fieldVals[i], fieldVals[i+1]
			if err := b.Put(h.fieldKey(field), val); err != nil {
				return err
			}
		}
		return tx.Bucket(h.bucket.bucketName).Put(h.rawKey(), nil)
	})
}

func (h *Hash) Remove(fields ...[]byte) error {
	return h.bucket.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(h.bucket.bucketName)
		for _, field := range fields {
			if err := b.Delete(h.fieldKey(field)); err != nil {
				return err
			}
		}
		// check if removeAll or not
		prefix := h.fieldPrefix()
		k, _ := b.Cursor().Seek(prefix)
		if !bytes.HasPrefix(k, prefix) {
			return b.Delete(h.rawKey())
		} else {
			return nil
		}
	})
}

func (h *Hash) Drop() error {
	return h.bucket.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(h.bucket.bucketName)
		c := b.Cursor()
		prefix := h.fieldPrefix()
		for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return b.Delete(h.rawKey())
	})
}

// +key,h
func (h *Hash) rawKey() []byte {
	return rawKey(h.key, HASH)
}

// h[key]field
func (h *Hash) fieldKey(field []byte) []byte {
	return bytes.Join([][]byte{h.fieldPrefix(), field}, nil)
}

// h[key]
func (h *Hash) fieldPrefix() []byte {
	return bytes.Join([][]byte{[]byte{byte(HASH)}, SOK, h.key, EOK}, nil)
}

// split h[key]field into field
func (h *Hash) fieldInKey(fieldKey []byte) []byte {
	right := bytes.Index(fieldKey, EOK)
	return fieldKey[right+1:]
}
