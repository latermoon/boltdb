package db

import (
	"github.com/boltdb/bolt"
)

func (b *BoltService) AllKeyValues() (map[string]string, error) {
	result := make(map[string]string)
	err := b.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(b.bucket).ForEach(func(k, v []byte) error {
			result[string(k)] = string(v)
			return nil
		})
	})
	return result, err
}
