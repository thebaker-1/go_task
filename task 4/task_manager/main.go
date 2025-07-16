package main

import (
	"fmt"
	"task_manager/router"
)

func main() {
	fmt.Println("Task Manager API")
	r := router.SetupRouter()
	r.Run()
}
