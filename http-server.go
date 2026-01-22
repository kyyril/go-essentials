package main

import (
	"fmt"
	"net/http"
	"time"
)

// Handler function for the root path
func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Hello, World!")
}

// Handler function for the greet path
func greetHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "Guest"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}

// Handler function for the headers path
func headersHandler(w http.ResponseWriter, r *http.Request) {
    for name, headers := range r.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

// stration function
func strateHTTPServer() {
    fmt.Println("=== HTTP Server Example ===")

    // Register handlers
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/greet", greetHandler)
    http.HandleFunc("/headers", headersHandler)

    // Start the server in a goroutine
    go func() {
        fmt.Println("Starting server on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            fmt.Printf("Server error: %v\n", err)
        }
    }()

    // Give the server a moment to start
    time.Sleep(100 * time.Millisecond)

    fmt.Println("Server is running. Try accessing:")
    fmt.Println("  http://localhost:8080/")
    fmt.Println("  http://localhost:8080/greet?name=Khairil")
    fmt.Println("  http://localhost:8080/headers")

    // Keep the program running
    select {}
}