package main

import (
	"curl/internal/url"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	elog = log.New(os.Stderr, "curl: ", 0)
	ilog = log.New(io.Discard, "", 0)

	fverbose = pflag.BoolP("verbose", "v", false, "Make the operation more talkative")
)

func main() {
	pflag.Parse()

	if *fverbose {
		ilog.SetOutput(os.Stderr)
	}

	run()
}

// GET /get HTTP/1.1
// Host: eu.httpbin.org
// Accept: */*
// Connection: close

var (
	defaultHdrs = map[string]string{
		"Accept": "*/*",
	}
)

func run() int {
	rawURL := pflag.Arg(0)
	if rawURL == "" {
		elog.Fatalf("try 'curl --help' for more information")
	}

	// parse URL
	url, err := url.Parse(rawURL)
	if err != nil {
		elog.Fatalf("parse error: %v", err)
	}

	// construct request
	reqbuf := new(strings.Builder)
	fmt.Fprintf(reqbuf, "GET %s HTTP/1.1\r\n", url.Path)
	fmt.Fprintf(reqbuf, "Host: %s\r\n", url.Host)
	for k, v := range defaultHdrs {
		fmt.Fprintf(reqbuf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(reqbuf, "Connection: close\r\n")
	reqbuf.WriteString("\r\n")

	// send request
	ilog.Printf("> %s", reqbuf.String())
	conn, err := net.Dial("tcp", url.Host+":"+url.PortOrDefault())
	if err != nil {
		elog.Fatalf("failed to establish connection: %v", err)
	}

	defer conn.Close()

	read, err := conn.Write([]byte(reqbuf.String()))
	if err != nil {
		elog.Fatalf("failed to send request: %v", err)
	}

	ilog.Printf("* sent %d bytes", read)

	respbuf := new(strings.Builder)
	wrote, err := io.Copy(respbuf, conn)
	if err != nil {
		elog.Fatalf("failed to read response: %v", err)
	}
	ilog.Printf("* received %d bytes", wrote)

	// dump response
	fmt.Print(respbuf.String())

	return 0
}
