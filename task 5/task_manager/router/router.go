package router

import (
	"github.com/gin-gonic/gin"
	"task_mdb/task_manager/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", controllers.GetTasks)
	r.GET("/tasks/:id", controllers.GetTask)
	r.PUT("/tasks/:id", controllers.UpdateTask)
	r.DELETE("/tasks/:id", controllers.DeleteTask)
	r.POST("/tasks", controllers.AddTask)

	return r
}
