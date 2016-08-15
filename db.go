package bolt

import (
	"sync"

	"github.com/boltdb/bolt"
)

// New ...
// db := bolt.New("/tmp/demo.db")
// db.Bucket("0").Hash("user:100422").Get("name")
// db.Bucket("0").List("acl:group:list").Range(0, 10)
func New(dbpath string) (*DB, error) {
	db, err := bolt.Open(dbpath, 0644, nil)
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
