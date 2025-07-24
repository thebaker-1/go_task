package router

import (
	"github.com/gin-gonic/gin"
	"task_mdb/controllers"
	"task_mdb/middleware"
)

func SetupRouter(ctrl *controllers.Controller) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", ctrl.RegisterUser)
	r.POST("/login", ctrl.LoginUser)

	// Protected routes (require authentication)
	protected := r.Group("/")
	protected.Use(middleware.AuthenticateJWT())
	{
		// Task management routes
		protected.GET("/tasks", controllers.GetTasks)
		protected.GET("/tasks/:id", controllers.GetTask)
		protected.POST("/tasks", controllers.AddTask)

		// Admin-only routes (require authentication and admin role)
		admin := protected.Group("/")
		admin.Use(middleware.AuthorizeRole("admin"))
		{
			admin.PUT("/tasks/:id", controllers.UpdateTask)
			admin.DELETE("/tasks/:id", controllers.DeleteTask)
		}
	}

	return r
}