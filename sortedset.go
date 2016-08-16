package boltdb

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

// SortedSet ...
// +key,z = ""
// z[key]m member = score
// z[key]s score member = ""
type SortedSet struct {
	bucket *Bucket
	key    []byte
}

// Add add score & member pairs
// SortedSet.Add(Score, []byte, Score, []byte ...)
func (s *SortedSet) Add(scoreMembers ...[]byte) (int, error) {
	count := len(scoreMembers)
	if count < 2 || count%2 != 0 {
		return 0, errors.New("invalid score/member pairs")
	}
	added := 0
	err := s.bucket.Update(func(b *bolt.Bucket) error {
		for i := 0; i < count; i += 2 {
			score, member := scoreMembers[i], scoreMembers[i+1]
			skey, mkey := s.scoreKey(score, member), s.memberKey(member)
			oldscore := b.Get(mkey)
			// remove old score key
			if oldscore != nil {
				oldskey := s.scoreKey(oldscore, member)
				if err := b.Delete(oldskey); err != nil {
					return err
				}
			} else {
				added++
			}
			b.Put(mkey, score)
			b.Put(skey, nil)
		}
		b.Put(s.rawKey(), nil)
		return nil
	})
	return added, err
}

func (s SortedSet) Score(member []byte) (Score, error) {
	var score []byte
	err := s.bucket.View(func(b *bolt.Bucket) error {
		score = b.Get(s.memberKey(member))
		return nil
	})
	return score, err
}

func (s *SortedSet) Remove(members ...[]byte) (int, error) {
	removed := 0 // not including non existing members
	err := s.bucket.Update(func(b *bolt.Bucket) error {
		for _, member := range members {
			score := b.Get(s.memberKey(member))
			if score == nil {
				continue
			}
			if err := b.Delete(s.scoreKey(score, member)); err != nil {
				return err
			}
			if err := b.Delete(s.memberKey(member)); err != nil {
				return err
			}
			removed++
		}
		// clean up
		prefix := s.keyPrefix()
		if k, _ := b.Cursor().Seek(prefix); !bytes.HasPrefix(k, prefix) {
			return b.Delete(s.rawKey())
		}
		return nil
	})
	return removed, err
}

// RevRangeByScore ...
// <fr> is larger than <to>
func (s *SortedSet) RevRangeByScore(fr, to Score, fn func(i int64, score Score, member []byte, quit *bool)) error {
	min := s.scorePrefix(to)
	max := append(s.scorePrefix(fr), MAXBYTE)
	return s.bucket.View(func(b *bolt.Bucket) error {
		c := b.Cursor()
		var i int64 // 0
		k, _ := c.Seek(max)
		for k, _ = c.Prev(); k != nil && bytes.Compare(k, min) >= 0; k, _ = c.Prev() {
			quit := false
			score, member, err := s.splitScoreKey(k)
			if err != nil {
				return err
			}
			if fn(i, score, member, &quit); quit {
				break
			}
			i++
		}
		return nil
	})
}

// RangeByScore ...
// <fr> is less than <to>
func (s *SortedSet) RangeByScore(fr, to Score, fn func(i int64, score Score, member []byte, quit *bool)) error {
	min := s.scorePrefix(fr)
	max := append(s.scorePrefix(to), MAXBYTE)
	return s.bucket.View(func(b *bolt.Bucket) error {
		c := b.Cursor()
		var i int64 // 0
		for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
			quit := false
			score, member, err := s.splitScoreKey(k)
			if err != nil {
				return err
			}
			if fn(i, score, member, &quit); quit {
				break
			}
			i++
		}
		return nil
	})
}

// +key,z = ""
func (s *SortedSet) rawKey() []byte {
	return rawKey(s.key, ElemType(SORTEDSET))
}

// z[key]
func (s *SortedSet) keyPrefix() []byte {
	return bytes.Join([][]byte{[]byte{byte(SORTEDSET)}, SOK, s.key, EOK}, nil)
}

// z[key]m
func (s *SortedSet) memberKey(member []byte) []byte {
	return bytes.Join([][]byte{s.keyPrefix(), []byte{'m'}, member}, nil)
}

// z[key]s score
func (s *SortedSet) scorePrefix(score []byte) []byte {
	return bytes.Join([][]byte{s.keyPrefix(), []byte{'s'}, score, []byte{' '}}, nil)
}

// z[key]s score member
func (s *SortedSet) scoreKey(score, member []byte) []byte {
	return bytes.Join([][]byte{s.keyPrefix(), []byte{'s'}, score, []byte{' '}, member}, nil)
}

// split (z[key]s score member) into (score, member)
func (s *SortedSet) splitScoreKey(skey []byte) ([]byte, []byte, error) {
	buf := bytes.TrimPrefix(skey, s.keyPrefix())
	pairs := bytes.Split(buf[1:], []byte{' '}) // skip score mark 's'
	if len(pairs) != 2 {
		return nil, nil, errors.New("invalid score/member key: " + string(skey))
	}
	return pairs[0], pairs[1], nil
}

// split (z[key]m member) into (member)
func (s *SortedSet) splitMemberKey(mkey []byte) ([]byte, error) {
	buf := bytes.TrimPrefix(mkey, s.keyPrefix())
	return buf[1:], nil // skip member mark 'm'
}
