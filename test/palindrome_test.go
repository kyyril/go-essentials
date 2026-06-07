package main

import (
	"testing"
)

func TestIsPalindrome (t *testing.T) {
	test := []struct{
		name string
		input string
		want bool
	}{
		{"valid palindrome", "katak", true},
		{"invalid palindrome", "golang", false},
		{"empty string", "", true},
		{"single character", "a", true},
	}
	for _, tc := range test {
		// t.Run creates a sub-test
		t.Run(tc.name, func(t *testing.T){
			got := IsPalindrome(tc.input)
			want := tc.want
			if got != want {
				t.Errorf("IsPalindrome(%q) = %t; want is %t", tc.input, got, want)
			}
		})
	}
}