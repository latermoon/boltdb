# go.bolt
go.roa wrapper for boltdb

# connect to go.blot
redis-cli -p 3002
> bolt.SelectBucket(name)
> bolt.Put(key, value)
> bolt.Get(key)
> bolt.Delete(key, ...)

import "github.com/latermoon/go.bolt"

db := bolt.New(dbpath)
db.Get("version")
db.Hash("user:100422:profile").HGet("name")
db.List("acl:group:rules").Range(0, 100)
db.SortedSet("userlist").ZAdd("score", "member", ...)

