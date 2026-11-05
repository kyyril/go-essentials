package main

import "testing"

func TestCalculatePrice(t *testing.T) {
	t.Run("should be return zero when price is zero", func(t *testing.T) {
		got := CalculatePrice(0)
		want := 0.0
		assertCorrectPrice(t, got, want)

	})
	t.Run("calculate price with tax", func(t *testing.T) {
		got := CalculatePrice(100)
		want := 110.0
		assertCorrectPrice(t, got, want)
	})
}

func assertCorrectPrice(tb testing.TB, got, want float64) {
	tb.Helper()
	if got != want {
		tb.Errorf("got %v want %v", got, want)
	}
}
