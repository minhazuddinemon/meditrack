package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	dbGen "meditrack/internal/db"
	"meditrack/internal/handler"
	"meditrack/templates/pages"
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

	mux := http.NewServeMux()

	// Serve static files (matches /static/*)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		pages.Landing().Render(r.Context(), w)
	})

	receptionHandler := &handler.ReceptionHandler{Queries: queries}
	mux.Handle("GET /reception", receptionHandler)
	mux.HandleFunc("POST /reception/token", receptionHandler.IssueToken)
	mux.HandleFunc("POST /reception/register-student", receptionHandler.RegisterStudent)

	// Doctor Module
	// Doctor Module
	doctorHandler := &handler.DoctorHandler{Queries: queries, RawDB: dbConn}
	mux.Handle("GET /doctor", doctorHandler)
	mux.HandleFunc("POST /doctor/login", doctorHandler.Login)
	mux.HandleFunc("GET /doctor/logout", doctorHandler.Logout)
	mux.HandleFunc("GET /doctor/patient/{st_id}", doctorHandler.LoadPatientWorkspace)
	mux.HandleFunc("GET /doctor/med-row", doctorHandler.GetMedicineRow)
	mux.HandleFunc("POST /doctor/prescription", doctorHandler.CreatePrescription)

	// Student Vault Module
	vaultHandler := &handler.VaultHandler{Queries: queries}
	mux.Handle("GET /vault", vaultHandler)
	mux.HandleFunc("GET /vault/history", vaultHandler.GetStudentHistory)
	mux.HandleFunc("GET /vault/rx/{p_id}", vaultHandler.ViewPrescription)

	// Root page ({$} ensures EXACT match on "/" so it doesn't overlap /static/)
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
