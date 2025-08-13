package main

import (
	"log"
	"net/http"
	"os"

	"taskmanager/Delivery/controllers"
	"taskmanager/Delivery/routers"
	Infrasturcture "taskmanager/Infrastructure"
	"taskmanager/Repositories"
	"taskmanager/Usecases"
)

func main() {
	log.Println("Starting Task Manager...")
	// Initialize MongoDB client using Infrastructure package
	Init()
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DATABASE_NAME")
	// log.Println("MongoDB URI:", mongoURI)
	tasksCollection := os.Getenv("TASKS_COLLECTION")
	usersCollection := os.Getenv("USERS_COLLECTION")
	mongoClient, err := Infrasturcture.NewMongoDBClient(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(); err != nil {
			log.Fatal("Failed to disconnect MongoDB:", err)
		}
	}()

	// Initialize repositories with collections from Infrastructure
	taskCollection := mongoClient.GetCollection(dbName, tasksCollection)
	userCollection := mongoClient.GetCollection(dbName, usersCollection)

	taskRepo := Repositories.NewMongoTaskRepository(&Repositories.MongoCollectionAdapter{Coll: taskCollection})
	userRepo := Repositories.NewMongoUserRepository(&Repositories.UserMongoCollectionAdapter{Coll: userCollection})

	// Initialize usecases
	taskUsecase := Usecases.NewTaskUsecase(taskRepo)
	userUsecase := Usecases.NewUserUsecase(userRepo)

	// Initialize controllers
	ctrl := controllers.NewController(userUsecase, taskUsecase)

	// Setup router
	r := routers.SetupRouter(ctrl)

	// Start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
