package main

import (
	"curl/internal/http"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
)

var (
	fverbose = pflag.BoolP("verbose", "v", false, "Make the operation more talkative")
)

func main() {
	pflag.Parse()
	log.SetFlags(0)
	if !*fverbose {
		log.SetOutput(io.Discard)
	}
	run()
}

func run() {
	rawURL := pflag.Arg(0)
	if rawURL == "" {
		log.Fatalf("try 'curl --help' for more information")
	}

	url, err := http.ParseURL(rawURL)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}

	req := http.NewRequest("GET", url, nil)

	resp, err := http.Do(req)
	if err != nil {
		log.Fatalf("http.Do error: %v", err)
	}

	fmt.Fprintf(os.Stdout, "%s", resp.Body)
}
