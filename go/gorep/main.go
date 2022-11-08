package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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
	pattern := flag.Arg(0)
	if *caseInsensitive {
		pattern = strings.ToLower(pattern)
	}

	scanner := bufio.NewScanner(reader)

	var lineCount int
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if *caseInsensitive {
			line = strings.ToLower(line)
		}

		if strings.Contains(line, pattern) {
			if *showLineCount {
				lineCount++
				continue
			}

			if *showLineNumbers {
				fmt.Println(lineNumber, line)
			} else {
				fmt.Println(line)
			}
		}
	}

	if *showLineCount {
		fmt.Println(lineCount)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error while scanning file", err)
		os.Exit(1)
	}
}
