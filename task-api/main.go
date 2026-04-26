package main

import (
	"log"
	"net/http"
	"os"
	"task-api/internal/database"
	"task-api/internal/handlers"
	"task-api/internal/middleware"
)

func main() {
	os.MkdirAll("./data", os.ModePerm)
	database.InitDB()

	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)

	// Protected Routes (Butuh Login)
	// Kita bungkus handler dengan middleware
	mux.Handle("/tasks", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetTasks), ""))
	mux.Handle("/tasks/create", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateTask), "user"))
	
	// Khusus Admin
	mux.Handle("/tasks/delete", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteTask), "admin"))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}