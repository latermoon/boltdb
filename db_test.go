package bolt

import (
	"github.com/facebookgo/ensure"
	"io/ioutil"
	// "log"
	"path"
	"testing"
)

func newBoltDB(t *testing.T) *DB {
	dir, err := ioutil.TempDir("", "bolt")
	ensure.Nil(t, err)

	dbpath := path.Join(dir, "bolt.db")
	// log.Println("dbpath:", dbpath)
	db, err := New(dbpath)
	ensure.Nil(t, err)

	return db
}
