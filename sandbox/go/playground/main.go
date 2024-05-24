package main

import "fmt"

func main() {
	ascii := [256]uint8{'a': 1, '\t': 1}
	fmt.Println(ascii['a'], ascii['\t'])
}
