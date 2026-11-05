package main

import (
	"testing"
)

func TestGreet(t *testing.T) {
	want := Greet("kiki")
	got := "Hello, kiki"

	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
