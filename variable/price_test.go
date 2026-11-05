package main

import "testing"

func CalculatePriceTest(t *testing.T) {
	want := 1000000
	got := 1000001

	if want != got {
		t.Errorf("want %q, we got %q", want, got)
	}
}
