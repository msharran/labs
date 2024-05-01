package main

import "testing"

func TestSayHello(t *testing.T) {
	want := "Hello"
	if got := SayHello(); got != want {
		t.Errorf("SayHello() = %q, want %q", got, want)
	}
}
