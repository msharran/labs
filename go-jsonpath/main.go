package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"k8s.io/client-go/util/jsonpath"
)

func main() {
	var path string
	var r io.Reader

	switch len(os.Args) {
	case 1:
		fatal(fmt.Errorf("no args given"))
	case 2:
		path = os.Args[1]
		r = os.Stdin
	default:
		path = os.Args[1]
		f, err := os.Open(os.Args[2])
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		r = f
	}

	b, err := io.ReadAll(r)
	if err != nil {
		fatal(err)
	}

	var v interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		fatal(err)
	}

	j := jsonpath.New("")
	err = j.Parse(path)
	if err != nil {
		fatal(err)
	}

	err = j.Execute(os.Stdout, v)
	if err != nil {
		fatal(err)
	}

	fmt.Println()
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}
