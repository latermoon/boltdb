package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/latermoon/boltdb"
	"github.com/latermoon/boltdb/cmd/redis-bolt/handler"
	redis "github.com/latermoon/go-redis-server"
)

// > redis-bolt -db /tmp/redis.bolt -p 6780
func main() {
	log.Println("redis-bolt")

	dir, _ := ioutil.TempDir("", "boltdb")
	db, err := boltdb.Open(path.Join(dir, "bolt.db"), 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	h := handler.New(db)

	config := redis.DefaultConfig().Port(6380).Handler(h)
	server, err := redis.NewServer(config)
	if err != nil {
		panic(err)
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
