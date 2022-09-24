package main

import (
	"github.com/msharran/labs/go-zip-unarchive/store"
	"github.com/msharran/labs/go-zip-unarchive/store/memory"
	"github.com/msharran/labs/go-zip-unarchive/store/s3"
)

func main() {
	mem := memory.NewStore()
	s3 := s3.NewStore()

	save("hello", []byte("world"), mem, s3)
}

func save(key string, data []byte, writers ...store.Writer) {
	for _, w := range writers {
		w.Write(key, data)
	}
}
