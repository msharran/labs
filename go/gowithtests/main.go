package main

import "fmt"

const (
	LanguageFrench = "French"

	HelloPrefixEnglish = "Hello, "
	HelloPrefixFrench  = "Bonjour, "
)

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
