package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHello(t *testing.T) {
	t.Run("greet name when argument is passed", func(t *testing.T) {
		got := Hello("Sharran", "")
		want := "Hello, Sharran"

		assertString(t, want, got)
	})

	t.Run("greet world by default when argument is empty", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"

		assertString(t, want, got)
	})

	t.Run("greet in english if language is empty", func(t *testing.T) {
		got := Hello("Sharran", "")
		want := "Hello, Sharran"

		assertString(t, want, got)
	})

	t.Run("greet in english if language is unsupported", func(t *testing.T) {
		got := Hello("Sharran", "Spanish")
		want := "Hello, Sharran"

		assertString(t, want, got)
	})

	t.Run("greet in french", func(t *testing.T) {
		got := Hello("Sharran", LanguageFrench)
		want := "Bonjour, Sharran"

		assertString(t, want, got)
	})

	t.Run("greet world in french", func(t *testing.T) {
		got := Hello("", LanguageFrench)
		want := "Bonjour, World"

		assertString(t, want, got)
	})
}

func TestGreet(t *testing.T) {
	w := bytes.Buffer{}
	Greet(&w, "Sharran")
	want := "Hello, Sharran"
	got := strings.TrimSpace(w.String())

	assertString(t, want, got)
}

type delayedServer struct {
	delay time.Duration
}

func (s *delayedServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(s.delay)
}

func TestPingRace(t *testing.T) {
	t.Run("get fastest site", func(t *testing.T) {
		slowServer := httptest.NewServer(&delayedServer{delay: 10 * time.Millisecond})
		defer slowServer.Close()

		fastServer := httptest.NewServer(&delayedServer{delay: 5 * time.Millisecond})
		defer fastServer.Close()

		got, err := RacePing(slowServer.URL, fastServer.URL)

		assertErr(t, nil, err)
		assertString(t, fastServer.URL, got)
	})

	t.Run("fail if both requests timedout", func(t *testing.T) {
		slowServer := httptest.NewServer(&delayedServer{delay: 10 * time.Millisecond})
		defer slowServer.Close()

		fastServer := httptest.NewServer(&delayedServer{delay: 5 * time.Millisecond})
		defer fastServer.Close()

		_, err := TimedRacePing(slowServer.URL, fastServer.URL, 2*time.Millisecond)

		assertErr(t, ErrReqTimedOut, err)
	})
}

func assertErr(t testing.TB, want, got error) {
	t.Helper()
	if got != want {
		t.Errorf("wantErr: %v, gotErr: %v", want, got)
	}
}

func assertString(t testing.TB, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("len(want): %d, len(got): %d", len(want), len(got))
		t.Errorf("want: %q, got: %q", want, got)
	}
}
