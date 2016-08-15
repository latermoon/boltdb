package bolt

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

// List ...
// +key,l = ""
// l[key]0 = "a"
// l[key]1 = "b"
// l[key]2 = "c"
type List struct {
	bucket *Bucket
	key    []byte
}

// RPush ...
func (l *List) RPush(vals ...[]byte) error {
	x, y, err := l.rangeIndex()
	if err != nil {
		return err
	}
	err = l.bucket.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(l.bucket.bucketName)
		if x == 0 && y == -1 {
			b.Put(l.rawKey(), nil)
		}
		for i, val := range vals {
			b.Put(l.indexKey(y+int64(i)+1), val)
		}
		return nil
	})
	return err
}

// LPush ...
func (l *List) LPush(vals ...[]byte) error {
	x, y, err := l.rangeIndex()
	if err != nil {
		return err
	}
	err = l.bucket.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(l.bucket.bucketName)
		if x == 0 && y == -1 {
			b.Put(l.rawKey(), nil)
		}
		for i, val := range vals {
			b.Put(l.indexKey(x-int64(i)-1), val)
		}
		return nil
	})
	return err
}

// RPop ...
func (l *List) RPop() ([]byte, error) {
	return l.pop(false)
}

// LPop ...
func (l *List) LPop() ([]byte, error) {
	return l.pop(true)
}

func (l *List) pop(left bool) ([]byte, error) {
	x, y, err := l.rangeIndex()
	if err != nil {
		return nil, err
	}

	size := y - x + 1
	if size == 0 {
		return nil, nil
	} else if size < 0 { // double check
		return nil, errors.New("bad list struct")
	}

	var idxkey []byte
	if left {
		idxkey = l.indexKey(x)
	} else {
		idxkey = l.indexKey(y)
	}

	var val []byte
	err = l.bucket.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(l.bucket.bucketName)
		val = b.Get(idxkey)
		if err := b.Delete(idxkey); err != nil {
			return err
		}
		if size == 1 { // clean up
			return b.Delete(l.rawKey())
		}
		return nil
	})

	return val, nil
}

// Len ...
func (l *List) Len() (int64, error) {
	x, y, err := l.rangeIndex()
	return y - x + 1, err
}

func (l *List) rangeIndex() (int64, int64, error) {
	left, err := l.leftIndex()
	if err != nil {
		return 0, -1, err
	}
	right, err := l.rightIndex()
	if err != nil {
		return 0, -1, err
	}
	return left, right, nil
}

func (l *List) leftIndex() (int64, error) {
	idx := int64(0) // default 0
	err := l.bucket.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(l.bucket.bucketName).Cursor()
		prefix := l.keyPrefix()
		if k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix) {
			idx = l.indexInKey(k)
		}
		return nil
	})
	return idx, err
}

func (l *List) rightIndex() (int64, error) {
	idx := int64(-1) // default -1
	err := l.bucket.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(l.bucket.bucketName).Cursor()
		prefix := l.keyPrefix()
		seek := append(prefix, MAXBYTE)
		k, _ := c.Seek(seek)
		// move to previous record & check
		if k, _ = c.Prev(); bytes.HasPrefix(k, prefix) {
			idx = l.indexInKey(k)
		}
		return nil
	})
	return idx, err
}

// +key,l = ""
func (l *List) rawKey() []byte {
	return rawKey(l.key, ElemType(LIST))
}

// l[key]
func (l *List) keyPrefix() []byte {
	return bytes.Join([][]byte{[]byte{byte(LIST)}, SOK, l.key, EOK}, nil)
}

// l[key]0 = "a"
func (l *List) indexKey(i int64) []byte {
	sign := []byte{0}
	if i >= 0 {
		sign = []byte{1}
	}
	return bytes.Join([][]byte{l.keyPrefix(), sign, Int64ToBytes(i)}, nil)
}

// split l[key]index into index
func (l *List) indexInKey(key []byte) int64 {
	idxbuf := bytes.TrimPrefix(key, l.keyPrefix())
	return BytesToInt64(idxbuf[1:]) // skip sign "0/1"
}
