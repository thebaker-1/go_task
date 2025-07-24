package main

import (
	// "context"
	"fmt"
	"log"
	"task_mdb/controllers"
	"task_mdb/data"
	"task_mdb/router"
)

func main() {

	err := data.InitMongo("mongodb://localhost:27017/", "testdb", "tasks")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	// Initialize task service collection
	data.InitTaskService(data.Client, "testdb", "tasks")
	defer func() {
		if err := data.DisconnectMongo(); err != nil {
			log.Fatal("Failed to disconnect MongoDB:", err)
		}
	}()

	fmt.Println("Task Manager API")

	userService := data.NewUserService(data.Client, "testdb", "users")
	ctrl := controllers.NewController(userService)

	r := router.SetupRouter(ctrl)
	r.Run()
}