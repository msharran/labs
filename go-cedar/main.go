package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	// flag entities-json-file
	entities string
	// flag request-json-file
	request string
	// flag policy-cedar-file
	policies string
)

func init() {
	flag.StringVar(&entities, "entities", "", "entities json file")
	flag.StringVar(&request, "request-json", "", "request json file")
	flag.StringVar(&policies, "policies", "", "policy cedar file")
}

func main() {
	flag.Parse()

	cedarBin := os.Getenv("CEDAR_BIN")
	if cedarBin == "" {
		cedarBin = "cedar"
	}

	out, err := exec.Command(cedarBin, "authorize", "--entities", entities, "--request-json", request, "--policies", policies).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
