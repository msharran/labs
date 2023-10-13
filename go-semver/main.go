package main

import (
	"fmt"
	"os"

	"golang.org/x/mod/semver"
)

func main() {
	cmp := semver.Compare(os.Args[1], os.Args[2])
	fmt.Println(cmp)
}
