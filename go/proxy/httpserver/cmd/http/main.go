package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling request", r.Host)
		fmt.Fprintln(w, "Hello, World")
	})

	log.Printf("started the server at port %d\n", 8080)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
