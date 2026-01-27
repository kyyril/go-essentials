package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("number fibo:")
	for i := 1; i < 40; i++ {
		fmt.Println(fibonacci(i))
	}
	
	// Test performance with caching
	fmt.Println("Testing performance for fibonacci(100) x10000 (with cache):")
	start := time.Now()
	for i := 0; i < 10000; i++ {
		_ = fibonacci(100)
	}
	elapsed := time.Since(start)
	fmt.Println("Time for 10000 calls with cache:", elapsed.Nanoseconds(), "ns")
	fmt.Println("Cache size after:", len(fibCache))

	// Test performance without caching
	fmt.Println("Testing performance for fibonacciNoCache(100) x10000 (without cache):")
	start = time.Now()
	for i := 0; i < 10000; i++ {
		_ = fibonacciNoCache(100)
	}
	elapsed = time.Since(start)
	fmt.Println("Time for 10000 calls without cache:", elapsed.Nanoseconds(), "ns")
}

var fibCache = make(map[int]int)

func fibonacci(n int) int {
	if val, ok := fibCache[n]; ok {
		return val
	}
	if n == 0 {
		fibCache[0] = 0
		return 0
	}
	if n == 1 {
		fibCache[1] = 1
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	fibCache[n] = b
	return b
}

func fibonacciNoCache(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}