package main

import (
	"fmt"
	"io"
	"strings"
)

type mem struct {
	data []byte
}

func (s *mem) Write(p []byte) (n int, err error) {
	s.data = append(s.data, p...)
	return len(p), nil
}

func main() {
	m := &mem{}

	fmt.Println(string(m.data))

	_, err := io.Copy(m, strings.NewReader("foo "))
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(m, strings.NewReader("bar"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(m.data))
}
