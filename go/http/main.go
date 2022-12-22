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
	err := http.ListenAndServeTLS(":8080", "example.com+2.pem", "example.com+2-key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}
