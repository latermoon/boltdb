# go.bolt
go.roa wrapper for boltdb

import (
	"github.com/latermoon/go.bolt"
	"github.com/latermoon/go.bolt/score"
)

db, err := bolt.New(dbpath)
db.Set("version", "0.1.3")
db.Get("version")
db.Hash("user:100422:profile").Get("name")
db.List("acl:group:rules").RPush("a", "b", "c")
db.List("acl:group:rules").Range(0, 2)
db.SortedSet("userlist").Add("score", "member", ...)
db.TypeOf("key")

score := bolt.FloatScore(100.35)
score.Incr(100.35) 
zset.Add(score.Bytes(), "100422")

score := zset.Score("100422")
i := score.Int()



