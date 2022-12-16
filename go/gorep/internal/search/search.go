package search

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ExecArgs struct {
	Pattern         string
	CaseInsensitive bool
	ShowLineCount   bool
	ShowLineNumbers bool
}

func Exec(w io.Writer, r io.Reader, args ExecArgs) error {
	pattern := args.Pattern
	scanner := bufio.NewScanner(r)

	var lineCount int
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if args.CaseInsensitive {
			pattern = strings.ToLower(pattern)
			line = strings.ToLower(line)
		}

		if strings.Contains(line, pattern) {
			if args.ShowLineCount {
				lineCount++
				continue
			}

			if args.ShowLineNumbers {
				fmt.Fprintf(w, "%d %s\n", lineNumber, scanner.Text())
			} else {
				fmt.Fprintf(w, "%s\n", scanner.Text())
			}
		}
	}

	if args.ShowLineCount {
		fmt.Fprintln(w, lineCount)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
