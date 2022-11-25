package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	LanguageFrench = "French"

	HelloPrefixEnglish = "Hello, "
	HelloPrefixFrench  = "Bonjour, "
)

var ErrReqTimedOut = fmt.Errorf("request timed out")

func main() {
	fmt.Println(Hello("", ""))
}

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

func RacePing(a, b string) (string, error) {
	return TimedRacePing(a, b, 10*time.Second)
}

func TimedRacePing(a, b string, timeout time.Duration) (string, error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", ErrReqTimedOut
	}
}

func ping(url string) (res chan struct{}) {
	res = make(chan struct{}, 1)
	go func() {
		http.Get(url)
		close(res)
	}()
	return
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
