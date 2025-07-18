package main

import (
	"fmt"
	"log"
	"task_mdb/task_manager/data"
	"task_mdb/task_manager/router"
)

func main() {

	err := data.InitMongo("mongodb://localhost:27017/", "testdb", "tasks")
if err != nil {
    log.Fatal("Failed to connect to MongoDB:", err)
}
defer func() {
    if err := data.DisconnectMongo(); err != nil {
        log.Fatal("Failed to disconnect MongoDB:", err)
    }
}()

	fmt.Println("Task Manager API")
	r := router.SetupRouter()
	r.Run()
}
