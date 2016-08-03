package main

import (
	"github.com/latermoon/go.bolt/db"
	"github.com/latermoon/go.roa"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("no dbpath specified.")
		return
	}

	addr := os.Args[1]
	dbpath := os.Args[2]
	log.Println("go.blot listen to", addr, "dbpath:", dbpath)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	serv, err := db.New(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	roa.RegisterName("bolt", serv)
	roa.Serve(lis)
}
