# go.bolt
go.roa wrapper for boltdb

import "github.com/latermoon/go.bolt"

db := bolt.New(dbpath)
db.Get("version")
db.Hash("user:100422:profile").Get("name")
db.List("acl:group:rules").Range(0, 100)
db.SortedSet("userlist").Add("score", "member", ...)
db.TypeOf("version")


