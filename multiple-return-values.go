package main

import "fmt"

// Function that returns multiple values: quotient and error
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Function that returns multiple values: sum and product
func calculate(a, b int) (int, int) {
    return a + b, a * b
}

// Example function strating multiple return values
func strateMultipleReturns() {
    // Example 1: Using divide function
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Division result:", result)
    }

    // Example 2: Using calculate function
    sum, product := calculate(3, 4)
    fmt.Printf("Sum: %d, Product: %d\n", sum, product)

    // Example 3: Ignoring one return value with blank identifier
    _, prod := calculate(5, 6)
    fmt.Println("Product only:", prod)
}