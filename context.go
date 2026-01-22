package main

import (
	"context"
	"fmt"
	"time"
)

// Function strating context with timeout
func fetchData(ctx context.Context, id int) {
    select {
    case <-time.After(2 * time.Second):
        fmt.Printf("Fetched data for ID %d\n", id)
    case <-ctx.Done():
        fmt.Printf("Request for ID %d cancelled: %v\n", id, ctx.Err())
    }
}

// Function strating context with value
func processRequest(ctx context.Context) {
    value := ctx.Value("userID")
    if value != nil {
        fmt.Printf("Processing request for user: %v\n", value)
    } else {
        fmt.Println("No userID in context")
    }
}

// stration function
func strateContext() {
    fmt.Println("=== Context Examples ===")

    // Example 1: Context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    go fetchData(ctx, 1)
    go fetchData(ctx, 2)

    // Give time for goroutines to complete
    time.Sleep(3 * time.Second)

    // Example 2: Context with value
    ctx2 := context.WithValue(context.Background(), "userID", 12345)
    processRequest(ctx2)

    // Example 3: Context with cancellation
    ctx3, cancel3 := context.WithCancel(context.Background())
    go func() {
        time.Sleep(1 * time.Second)
        cancel3()
    }()

    select {
    case <-ctx3.Done():
        fmt.Printf("Context cancelled: %v\n", ctx3.Err())
    }
}