package main

import (
	"fmt"
)

func main() {
	// Ref "strings.Fields" for
	// ascii ([256]uint8) example
	// ascii := [256]uint8{'a': 1, '\t': 1}
	// fmt.Println(ascii['a'], ascii['\t'])

	// https://go.dev/blog/strings
	var a int32 = 0x2318
	fmt.Printf("%d\n", a)
	fmt.Printf("%c\n", rune(a))
	fmt.Printf("%q\n", rune(a))
	fmt.Printf("%+q\n", rune(a))
	fmt.Printf("%U\n", rune(a))
	fmt.Printf("%x\n", rune(a))
	for i := 0; i < len(`⌘`); i++ {
		fmt.Printf("%x ", `⌘`[i])
	}
	fmt.Println()

	uni := `Emojis: ⌘✅☕`
	fmt.Printf("%s\n", uni)

	// print the UTF-8 encoded bytes
	// of the string literal, not the
	// characters
	for i := 0; i < len(uni); i++ {
		fmt.Printf("%x ", uni[i])
	}
	fmt.Println()

	// decode unicode codepoints in
	// a string literal
	// i gives the byte position of the start of the rune.
	// runeValue gives the codepoint of the rune.
	for i, runeValue := range uni {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
	}
	fmt.Println()
}
