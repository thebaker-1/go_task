package main

// package main

import (
	"log"
	"os"
	"taskmanager/Infrastructure"
	"github.com/joho/godotenv"
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
