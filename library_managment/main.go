package main

import (
	"library_managment/controllers"
	"library_managment/models"
	"library_managment/services"
)

func main() {
	// Initialize the library service
	libraryService := services.NewLibrary()

	// Add some initial members (not through console, just for demonstration)
	libraryService.AddMember(models.Member{Name: "Alice Smith"})
	libraryService.AddMember(models.Member{Name: "Bob Johnson"})

	// Create the controller and run the application
	controller := controllers.NewLibraryController(libraryService)
	controller.Run()
}