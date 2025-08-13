package routers

import (
	"taskmanager/Delivery/controllers"
	"taskmanager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *controllers.Controller) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", ctrl.RegisterUser)
	r.POST("/login", ctrl.LoginUser)

	// Protected routes
	auth := r.Group("/")
	auth.Use(Infrastructure.AuthenticateJWT())

	auth.GET("/tasks", ctrl.GetTasks)
	auth.GET("/tasks/:id", ctrl.GetTask)
	auth.POST("/tasks", ctrl.AddTask)
	auth.PUT("/tasks/:id", ctrl.UpdateTask)
	auth.DELETE("/tasks/:id", ctrl.DeleteTask)

	return r
}
