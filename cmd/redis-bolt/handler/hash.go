package handler

func (r *RedisHandler) HSET(key, field, value []byte) error {
	hash, err := r.bucket().Hash(key)
	if err != nil {
		return err
	}
	return hash.Set(field, value)
}

func (r *RedisHandler) HGET(key, field []byte) ([]byte, error) {
	hash, err := r.bucket().Hash(key)
	if err != nil {
		return nil, err
	}
	return hash.Get(field)
}

func (r *RedisHandler) HGETALL(key []byte) ([][]byte, error) {
	hash, err := r.bucket().Hash(key)
	if err != nil {
		return nil, err
	}
	dict, err := hash.GetAll()
	if err != nil {
		return nil, err
	}
	var bulks [][]byte
	for k, v := range dict {
		bulks = append(bulks, []byte(k), v)
	}
	return bulks, nil
}
