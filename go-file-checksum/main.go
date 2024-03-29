package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: hash <file>")
	}

	file := os.Args[1]
	if _, err := os.Stat(file); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x", h.Sum(nil))
}
