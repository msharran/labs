package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var jsonpath string

func init() {
	flag.StringVar(&jsonpath, "jsonpath", "", "path to json file")
}

func main() {
	// start a http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	fmt.Println("Server started at port 8080")

	dir := os.Getenv("APP_DIR")
	if dir == "" {
		dir = "."
	}

	// create a tempfile
	f, err := os.CreateTemp(dir, "example")
	check(err)
	defer os.Remove(f.Name())

	fmt.Println("Temp file name:", f.Name())

	http.ListenAndServe(":8080", nil)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
