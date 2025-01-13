package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	inclHeaders = false
)

func init() {
	flag.BoolVar(&inclHeaders, "i", false, "Include headers in output")
}

var methods = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
	"TRACE":   true,
	"get":     true,
	"head":    true,
	"post":    true,
	"put":     true,
	"delete":  true,
	"connect": true,
	"options": true,
	"trace":   true,
}

var method = "GET"

func main() {
	flag.Parse()
	// get arg1 as argument
	arg1 := flag.Arg(0)
	if arg1 == "" {
		fatalf("url is required\n")
	}

	// if arg1 is a valid method, use it as method
	// and mandate url as arg2
	if methods[arg1] {
		method = strings.ToUpper(arg1)
		arg1 = flag.Arg(1)
		if arg1 == "" {
			fatalf("url is required\n")
		}
	}

	// create new request
	req, err := http.NewRequest(method, arg1, nil)
	if err != nil {
		fatalf("failed to create request: %v\n", err)
	}

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fatalf("failed to send request: %v\n", err)
	}
	defer resp.Body.Close()

	// print response
	if inclHeaders {
		fmt.Printf("%v %s\n", resp.Proto, resp.Status)
		for k, vv := range resp.Header {
			for _, v := range vv {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	}

	// copy response body to stdout
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fatalf("failed to copy response body: %v\n", err)
	}
	fmt.Println()
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	flag.Usage()
	os.Exit(1)
}
