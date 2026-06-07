package main

import (
	"testing"
)

func TestIsPalindrome (t *testing.T) {
	got := IsPalindrome("katak")
	want := true
	if got != want {
		t.Errorf("IsPalindrome('katak') = %t; want %t", got, want)
	}

	got = IsPalindrome("golang")
	want = false
	if got != want {
		t.Errorf("IsPalindrome('golang') = %t; want %t", got, want)
	}
}