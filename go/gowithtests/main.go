package main

import "fmt"

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
