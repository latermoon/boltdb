package bolt

import (
	"os"
	"sync"

	"github.com/boltdb/bolt"
)

// Options is another name of bolt.Options
// Options {Timeout: 0, ReadOnly: false}
type Options bolt.Options

// Open ...
// <mode> 0644 means -rw-r--r--
// <opt> alias of bolt.Options
// db, err := bolt.Open("/tmp/demo.db", 0644, nil)
// bucket, err := db.Bucket("0")
//
// hash, err := bucket.Hash("user:100422")
// val, err := hash.Get("name")
// list, err := bucket.List("userlist")
// list.RPush("a", "b", "c")
// item, err := list.LPop()
func Open(dbpath string, mode os.FileMode, opt *Options) (*DB, error) {
	db, err := bolt.Open(dbpath, mode, (*bolt.Options)(opt))
	if err != nil {
		return nil, err
	}
	return &DB{
		db:      db,
		buckets: map[string]*Bucket{},
	}, nil
}

type DB struct {
	db      *bolt.DB
	mu      sync.Mutex
	buckets map[string]*Bucket
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Bucket(name []byte) (*Bucket, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	bucket, exists := d.buckets[string(name)]
	if !exists {
		err := d.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(name)
			return err
		})
		if err != nil {
			return nil, err
		}
		bucket = &Bucket{db: d.db, bucketName: name}
		d.buckets[string(name)] = bucket
	}

	return bucket, nil
}
