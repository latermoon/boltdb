package handler

import "github.com/latermoon/boltdb"

// RedisHandler ...
type RedisHandler struct {
	db *boltdb.DB
}

func New(db *boltdb.DB) *RedisHandler {
	return &RedisHandler{db: db}
}

func (r *RedisHandler) bucket() *boltdb.Bucket {
	b, _ := r.db.Bucket([]byte("0"))
	return b
}

func (r *RedisHandler) GET(key []byte) ([]byte, error) {
	return r.bucket().Get(key)
}

func (r *RedisHandler) SET(key, value []byte) error {
	return r.bucket().Set(key, value)
}

func (r *RedisHandler) PING() (string, error) {
	return "PONG", nil
}
