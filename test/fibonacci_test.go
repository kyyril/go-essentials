package main

import (
	"testing"
)

func Test(t *testing.T) {
	actual:= fibonacci(4)
	
	if actual != 3 {
		t.Error("expection value 3 at position: 4")
	}
}