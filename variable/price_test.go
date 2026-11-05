package main

import "testing"

func TestCalculatePrice(t *testing.T) {
	t.Run("should be return zero when price is zero", func(t *testing.T) {
		got := CalculatePrice(0)
		want := 0.0

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})
	t.Run("calculate price with tax", func(t *testing.T) {
		got := CalculatePrice(100)
		want := 110.0

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}
