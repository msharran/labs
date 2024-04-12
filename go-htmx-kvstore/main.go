package main

import (
	"flag"
	"fmt"
	"go-htmx-kvstore/internal/server"
	"log"
)

var (
	port      = 1323
	db        = "kvstore.sqlite3"
	localhost = true
)

func init() {
	flag.IntVar(&port, "svr.port", 1323, "port to listen on")
	flag.StringVar(&db, "db.file_name", "tmp/kvstore.sqlite3", "database file")
}

// https://go.dev/doc/articles/wiki/
func main() {
	flag.Parse()

	s, err := server.New(
		server.WithPort(fmt.Sprint(port)),
		server.WithDBFileName(db),
		server.WithLocalHost(),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.ListenAndServe())
}
