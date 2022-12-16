package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/msharran/labs/go/gorep/internal/search"
)

var (
	showLineCount   = flag.Bool("c", false, "Prints the searched line count. This will not print the search result")
	showLineNumbers = flag.Bool("l", false, "Prints line numbers along with the search result")
	caseInsensitive = flag.Bool("i", false, "Case insensitive search")
)

func main() {
	flag.Parse()

	var reader io.Reader

	switch flag.NArg() {
	case 0:
		fmt.Fprintln(os.Stderr, "mandatory argument <PATTERN> is missing")
		os.Exit(1)
	case 1:
		reader = os.Stdin
	default:
		f, err := os.Open(flag.Arg(1))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while opening file", flag.Arg(1), err)
			os.Exit(1)
		}
		reader = f
	}

	err := search.Exec(os.Stdout, reader, search.ExecArgs{
		Pattern:         flag.Arg(0),
		CaseInsensitive: *caseInsensitive,
		ShowLineCount:   *showLineCount,
		ShowLineNumbers: *showLineNumbers,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Search error", err)
		os.Exit(1)
	}
}
