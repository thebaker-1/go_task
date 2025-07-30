package main

// package main

import (
	"log"
	"os"
	"taskmanager/Infrasturcture"

	"github.com/joho/godotenv"
)

// ...existing code...
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded, environment variables may be missing.")
	}
	loadedSecret := os.Getenv("JWT_SECRET")
	Infrasturcture.SetJWTSecret(loadedSecret)
}
