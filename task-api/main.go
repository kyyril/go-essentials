package main

import (
	"log"
	"net/http"
	"os"
	"task-api/internal/database"
	"task-api/internal/handlers"
)
func main (){
	os.MkadirAll("./data", os.ModePerm)

	database.InitDB()
	mux := http.NewServeMux()

	mux.HandlerFunc("/register", handlers.Register)
	mux.HandlerFunc("/login", handlers.Login)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
