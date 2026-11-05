package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	t.Run("saying Hello to people", func(t *testing.T) {
		got := Hello("kiki")
		want := "Hello, kiki"
		assertCorrectMessage(t, got, want)

	})
	t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(tb testing.TB, got, want string) {
	tb.Helper()
	if got != want {
		tb.Errorf("got %q want %q", got, want)
	}
}
