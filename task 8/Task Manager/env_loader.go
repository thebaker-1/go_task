package main

// package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"taskmanager/Infrastructure"
)

// ...existing code...
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded, environment variables may be missing.")
	}
	// *** CRITICAL DEBUGGING LINE ADDED HERE ***
	loadedSecret := os.Getenv("JWT_SECRET")
	Infrastructure.SetJWTSecret(loadedSecret)
}
