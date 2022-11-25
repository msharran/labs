package main

import (
	"fmt"
	"io"
)

const (
	LanguageFrench = "French"

	HelloPrefixEnglish = "Hello, "
	HelloPrefixFrench  = "Bonjour, "
)

// Hello is a function that greets the [name]
// in the provided [language].
// If the language is not supported/empty, it will
// greet in English
func Hello(name, language string) string {
	if name == "" {
		name = "World"
	}

	return fmt.Sprintf("%s%s", translatedPrefix(language), name)
}

func Greet(w io.Writer, name string) {
	fmt.Fprintln(w, fmt.Sprintf("%s%s", HelloPrefixEnglish, name))
}

func translatedPrefix(language string) (prefix string) {
	switch language {
	case LanguageFrench:
		prefix = HelloPrefixFrench
	default:
		prefix = HelloPrefixEnglish
	}
	return
}

func main() {
	fmt.Println(Hello("", ""))
}
