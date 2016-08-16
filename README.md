# go.boltdb

```
import (
	"github.com/latermoon/boltdb"
)

func main() {
	db, err := boltdb.Open("my.db", 0644, nil)
	defer db.Close()

	bucket, err := db.Bucket([]byte("0"))
	hash, err := bucket.Hash([]byte("hash"))
	list, err := bucket.List([]byte("list"))
	zset, err := bucket.SortedSet([]byte("zset"))

	bucket.Set([]byte("key"), []byte("value"))
	hash.Set([]byte("field"), []byte("value"))
	list.RPush([]byte("a"), []byte("b"), []byte("c"))
	zset.Add(Int64ToScore(-1), []byte("a"), Int64ToScore(0), []byte("b"), Int64ToScore(1), []byte("c"))
	zset.Add(Float64ToScore(-1.5), []byte("a"), Float64ToScore(0f), []byte("b"), Float64ToScore(1.5), []byte("c"))
}

```



