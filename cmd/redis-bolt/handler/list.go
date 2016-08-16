package handler

func (r *RedisHandler) LPUSH(key []byte, values ...[]byte) error {
	list, err := r.bucket().List(key)
	if err != nil {
		return err
	}
	return list.LPush(values...)
}

func (r *RedisHandler) RPUSH(key []byte, values ...[]byte) error {
	list, err := r.bucket().List(key)
	if err != nil {
		return err
	}
	return list.RPush(values...)
}

func (r *RedisHandler) LPOP(key []byte) ([]byte, error) {
	list, err := r.bucket().List(key)
	if err != nil {
		return nil, err
	}
	return list.LPop()
}

func (r *RedisHandler) RPOP(key []byte) ([]byte, error) {
	list, err := r.bucket().List(key)
	if err != nil {
		return nil, err
	}
	return list.RPop()
}

func (r *RedisHandler) LRANGE(key []byte, start, stop int) ([][]byte, error) {
	list, err := r.bucket().List(key)
	if err != nil {
		return nil, err
	}
	var bulks [][]byte
	err = list.Range(int64(start), int64(stop), func(i int64, value []byte, quit *bool) {
		bulks = append(bulks, value)
	})
	return bulks, err
}
