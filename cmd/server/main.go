package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	dbGen "meditrack/internal/db"
	"meditrack/internal/handler"
)

func main() {
	cfg := dbGen.Config{
		User:     getEnv("DB_USER", "minhaz"),
		Password: getEnv("DB_PASS", "405060"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		DBName:   getEnv("DB_NAME", "meditrack"),
	}

	dbConn, err := dbGen.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MariaDB: %v", err)
	}
	defer dbConn.Close()

	log.Println("Successfully connected to MariaDB!")

	queries := dbGen.New(dbConn)

	// Initialize Token Handler
	tokenHandler := handler.NewTokenHandler(queries, dbConn)

	mux := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Root page ({$} ensures EXACT match on "/" so it doesn't overlap /static/)
	mux.HandleFunc("GET /{$}", tokenHandler.ShowKiosk)

	// API / HTMX Routes
	mux.HandleFunc("POST /tokens", tokenHandler.IssueToken)
	port := getEnv("PORT", "8080")
	fmt.Printf("Server listening on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
