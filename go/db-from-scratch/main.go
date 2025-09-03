package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		var input string
		fmt.Printf("dblite> ")
		fmt.Scanln(&input)

		if strings.HasPrefix(input, ".") {
			if err := runMetaCommand(input); err != nil {
				fmt.Println(err)
			}
			continue
		}

		statement, err := prepareStatement(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := executeStatement(statement); err != nil {
			fmt.Println(err)
		}
	}
}

func runMetaCommand(input string) error {
	switch input {
	case ".exit":
		os.Exit(0)
		return nil
	default:
		return fmt.Errorf("unrecognized command '%s'", input)
	}
}

type StatementType int

const (
	STATEMENT_INSERT StatementType = iota
	STATEMENT_SELECT
)

type Statement struct {
	Type        StatementType
	RowToInsert *Row
}

type Row struct {
	Id    int
	Name  string
	Email string
}

var (
	ErrPrepareSyntaxError    = fmt.Errorf("prepare statement: syntax error")
	ErrorUnrecognizedKeyword = fmt.Errorf("prepare statement: unrecognized keyword")
)

func prepareStatement(input string) (*Statement, error) {
	if strings.HasPrefix(input, "insert") {
		var row Row
		nparsed, err := fmt.Sscanf(input, "insert %d %s %s", &row.Id, &row.Name, &row.Email)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrPrepareSyntaxError, input)
		}
		if nparsed < 3 {
			return nil, fmt.Errorf("%w: %s", ErrPrepareSyntaxError, input)
		}

		return &Statement{Type: STATEMENT_INSERT, RowToInsert: &row}, nil
	}

	if strings.HasPrefix(input, "select") {
		return &Statement{Type: STATEMENT_SELECT}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrorUnrecognizedKeyword, input)
}

func executeStatement(statement *Statement) error {
	switch statement.Type {
	case STATEMENT_INSERT:
		fmt.Println("This is where we would do an insert.")
	case STATEMENT_SELECT:
		fmt.Println("This is where we would do a select.")
	}
	return nil
}
