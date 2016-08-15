package score

// IntScore ...
type IntScore int64

func (s Score) From(b []byte) error {
	return nil
}

func (s Score) Bytes() []byte {
	return nil
}

func (s Score) Incr(n int64) {
}
