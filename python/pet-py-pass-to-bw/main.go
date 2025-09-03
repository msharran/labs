package main

import "fmt"


type Bitwarden struct {
	Items []Item `json:"items"`
}

type Item struct {
	Type       int        `json:"type"`
	Name       string     `json:"name"`
	Login      Login      `json:"login"`
}

func main() {
	fmt.Println("vim-go")
}
