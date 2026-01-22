package main

import (
	"fmt"
	"time"
)

// Function to be executed as a goroutine
func printNumbers() {
    for i := 1; i <= 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("Number: %d\n", i)
    }
}

// Function to be executed as a goroutine
func printLetters() {
    for i := 'a'; i <= 'e'; i++ {
        time.Sleep(150 * time.Millisecond)
        fmt.Printf("Letter: %c\n", i)
    }
}

// Function strating channels
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d started job %d\n", id, j)
        time.Sleep(time.Second)
        fmt.Printf("Worker %d finished job %d\n", id, j)
        results <- j * 2
    }
}

// stration function
func strateGoroutines() {
    fmt.Println("=== Goroutines Examples ===")

    // Example 1: Basic goroutines
    go printNumbers()
    go printLetters()

    // Give goroutines time to complete
    time.Sleep(2 * time.Second)
    fmt.Println("Done with basic goroutines")

    // Example 2: Goroutines with channels
    const numJobs = 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // Start workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for a := 1; a <= numJobs; a++ {
        fmt.Printf("Result: %d\n", <-results)
    }

    fmt.Println("Done with channel example")
}