package main

import (
	"fmt"
	"os/exec"
	"testing"
)

// Example test for the main function
func Example_curl_request(t *testing.T) {
	cmd := exec.Command("curl", "--http2", "http://localhost:3000/hello/sharran", "-i")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error: %v\nOutput: %s", err, out)
	}
	fmt.Printf("%s", out)
	// Output:
	// HTTP/1.1 101 Switching Protocols
	// Connection: Upgrade
	// Upgrade: h2c
	//
	// HTTP/2 200
	// content-type: text/plain; charset=utf-8
	// content-length: 15
	// date: Fri, 16 Feb 2024 12:51:25 GMT
	//
	// Hello, sharran!%
}
