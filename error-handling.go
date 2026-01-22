package main

import (
	"errors"
	"fmt"
)

// Function that returns a custom error
func validateInput(input string) error {
    if input == "" {
        return errors.New("input cannot be empty")
    }
    return nil
}

// Function using fmt.Errorf
func processData(data int) error {
    if data < 0 {
        return fmt.Errorf("invalid data: %d, must be non-negative", data)
    }
    return nil
}

// Function strating panic and recover
func safeDivision(a, b float64) (result float64, err error) {
    defer func() {
        if r := recover(); 
		r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    if b == 0 {
        panic("division by zero")
    }
    result = a / b
    return
}

// stration function
func strateErrorHandling() {
    fmt.Println("=== Error Handling Examples ===")

    // Example 1: Basic error checking
    err := validateInput("")
    if err != nil {
        fmt.Println("Validation error:", err)
    }

    err = validateInput("valid")
    if err != nil {
        fmt.Println("Validation error:", err)
    } else {
        fmt.Println("Input is valid")
    }

    // Example 2: Formatted error
    err = processData(-5)
    if err != nil {
        fmt.Println("Processing error:", err)
    }

    // Example 3: Panic and recover
    result, err := safeDivision(10, 0)
    if err != nil {
        fmt.Println("Safe division error:", err)
    } else {
        fmt.Printf("Result: %.2f\n", result)
    }

    result, err = safeDivision(10, 2)
    if err != nil {
        fmt.Println("Safe division error:", err)
    } else {
        fmt.Printf("Result: %.2f\n", result)
    }
}
