package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/masudcsesust04/golang-jwt-auth/internal/db"
	"github.com/masudcsesust04/golang-jwt-auth/internal/handlers"
	"github.com/masudcsesust04/golang-jwt-auth/internal/utils"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("JWT_SECRET")
	databaseURL := os.Getenv("DATABASE_URL")
	fmt.Println("JWT SECRET Key:", apiKey)
	fmt.Println("Database URL:", databaseURL)

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	dbConn, err := db.NewDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}
	defer dbConn.Close()

	// Initialize user handler
	userHandler := handlers.NewUserHandler(dbConn)

	// Setup router
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/logout", utils.JWTMiddleware(userHandler.Logout)).Methods("POST")
	router.HandleFunc("/token/refresh", userHandler.RefreshToken).Methods("POST")

	// user routes
	router.HandleFunc("/users", utils.JWTMiddleware(userHandler.GetUsers)).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUsers).Methods("POST")
	router.HandleFunc("/users/{id}", utils.JWTMiddleware(userHandler.GetUser)).Methods("GET")
	router.HandleFunc("/users/{id}", utils.JWTMiddleware(userHandler.UpdateUser)).Methods("PUT")
	router.HandleFunc("/users/{id}", utils.JWTMiddleware(userHandler.DeleteUser)).Methods("DELETE")

	// Start server
	addr := ":8080"
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
