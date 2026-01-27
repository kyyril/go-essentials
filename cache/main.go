package main

import (
	"fmt"
	"time"
)

var fibCache = make(map[int]int)

func main() {
	n := 40

	// ===== WITH CACHE =====
	fmt.Println("Testing fibonacciCache(", n, ") x1000 (WITH cache)")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_ = fibonacciCache(n)
	}
	elapsed := time.Since(start)
	fmt.Println("Time with cache:", elapsed)
	fmt.Println("Cache size:", len(fibCache))

	// reset cache
	fibCache = make(map[int]int)

	// ===== WITHOUT CACHE =====
	fmt.Println("\nTesting fibonacciNoCache(", n, ") x10 (WITHOUT cache)")
	start = time.Now()
	for i := 0; i < 10; i++ { // just 10 :(
		_ = fibonacciNoCache(n)
	}
	elapsed = time.Since(start)
	fmt.Println("Time without cache:", elapsed)
}

// ===== RECURSIVE + CACHE (MEMOIZATION) =====
func fibonacciCache(n int) int {
	if v, ok := fibCache[n]; ok {
		return v
	}

	var result int
	if n < 2 {
		result = n
	} else {
		result = fibonacciCache(n-1) + fibonacciCache(n-2)
	}

	fibCache[n] = result
	return result
}

// ===== RECURSIVE PURE (NO CACHE) =====
func fibonacciNoCache(n int) int {
	if n < 2 {
		return n
	}
	return fibonacciNoCache(n-1) + fibonacciNoCache(n-2)
}
