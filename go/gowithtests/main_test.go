package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	t.Run("greet name when argument is passed", func(t *testing.T) {
		got := Hello("Sharran", "")
		want := "Hello, Sharran"

		requireString(t, want, got)
	})

	t.Run("greet world by default when argument is empty", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"

		requireString(t, want, got)
	})

	t.Run("greet in english if language is empty", func(t *testing.T) {
		got := Hello("Sharran", "")
		want := "Hello, Sharran"

		requireString(t, want, got)
	})

	t.Run("greet in english if language is unsupported", func(t *testing.T) {
		got := Hello("Sharran", "Spanish")
		want := "Hello, Sharran"

		requireString(t, want, got)
	})

	t.Run("greet in french", func(t *testing.T) {
		got := Hello("Sharran", LanguageFrench)
		want := "Bonjour, Sharran"

		requireString(t, want, got)
	})

	t.Run("greet world in french", func(t *testing.T) {
		got := Hello("", LanguageFrench)
		want := "Bonjour, World"

		requireString(t, want, got)
	})
}

func TestGreet(t *testing.T) {
	w := bytes.Buffer{}
	Greet(&w, "Sharran")
	want := "Hello, Sharran"
	got := strings.TrimSpace(w.String())

	requireString(t, want, got)
}

func requireString(t testing.TB, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("len(want): %d, len(got): %d", len(want), len(got))
		t.Errorf("want: %s\ngot: %s", want, got)
	}
}
