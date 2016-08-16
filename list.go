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

func (l *List) Index(i int64) ([]byte, error) {
	x, err := l.leftIndex()
	if err != nil {
		return nil, err
	}
	var val []byte
	err = l.bucket.View(func(b *bolt.Bucket) error {
		val = b.Get(l.indexKey(x + i))
		return nil
	})
	return val, err
}

// Range enumerate value by index
// <start> must >= 0
// <stop> should equal to -1 or lager than <start>
func (l *List) Range(start, stop int64, fn func(i int64, value []byte, quit *bool)) error {
	if start < 0 || (stop != -1 && start > stop) {
		return errors.New("bad start/stop index")
	}
	x, y, err := l.rangeIndex()
	if err != nil {
		return err
	}
	if stop == -1 {
		stop = (y - x + 1) - 1 // (size) - 1
	}
	min := l.indexKey(x + int64(start))
	max := l.indexKey(x + int64(stop))
	return l.bucket.View(func(b *bolt.Bucket) error {
		c := b.Cursor()
		var i int64 // 0
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			quit := false
			if fn(start+i, v, &quit); quit {
				break
			}
			i++
		}
		return nil
	})
}

// RPush ...
func (l *List) RPush(vals ...[]byte) error {
	x, y, err := l.rangeIndex()
	if err != nil {
		return err
	}
	err = l.bucket.Batch(func(b *bolt.Bucket) error {
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
	err = l.bucket.Batch(func(b *bolt.Bucket) error {
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
	err = l.bucket.Batch(func(b *bolt.Bucket) error {
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
	err := l.bucket.View(func(b *bolt.Bucket) error {
		c := b.Cursor()
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
	err := l.bucket.View(func(b *bolt.Bucket) error {
		c := b.Cursor()
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
	return bytes.Join([][]byte{l.keyPrefix(), sign, itob(i)}, nil)
}

// split l[key]index into index
func (l *List) indexInKey(key []byte) int64 {
	idxbuf := bytes.TrimPrefix(key, l.keyPrefix())
	return btoi(idxbuf[1:]) // skip sign "0/1"
}
