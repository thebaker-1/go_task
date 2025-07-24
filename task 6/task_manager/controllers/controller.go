package controllers

import (
	"net/http"
	"time"

	"task_mdb/data"
	"task_mdb/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Define your JWT secret key. In a real app, this should be an environment variable.
var jwtSecret = []byte("your-super-secret-jwt-key")

type Controller struct {
	UserService *data.UserService
}

func NewController(userSvc *data.UserService) *Controller {
	return &Controller{UserService: userSvc}
}

func (c *Controller) RegisterUser(ctx *gin.Context) {
	var newUser models.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.UserService.RegisterUser(&newUser); err != nil {
		if err.Error() == "username already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (c *Controller) LoginUser(ctx *gin.Context) {
	var loginCredentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&loginCredentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserService.AuthenticateUser(loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": signedToken})
}

// GetTasks handles GET /tasks
func GetTasks(c *gin.Context) {
	tasks,err := data.GetAllTasks()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve tasks"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id
func GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	// Convert ObjectID to hex string for JSON response
	c.IndentedJSON(http.StatusOK, task)
}

// UpdateTask handles PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		Status      string `json:"status"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	task, err := data.GetTaskByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	updated := false

	if input.Title != "" {
		task.Title = input.Title
		updated = true
	}
	if input.Description != "" {
		task.Description = input.Description
		updated = true
	}
	if input.DueDate != "" {
		layout := "02-01-2006"
		parsedTime, err := time.Parse(layout, input.DueDate)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The time format is not 02-01-2006"})
			return
		}
		task.DueDate = parsedTime
		updated = true
	}
	if input.Status != "" {
		if input.Status == "Pending" || input.Status == "In Progress" || input.Status == "Completed" {
			task.Status = input.Status
			updated = true
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid status value"})
			return
		}
	}

	if !updated {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No valid fields to update"})
		return
	}

	updatedTask, err := data.UpdateTask(id, *task)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedTask)
}

// DeleteTask handles DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := data.DeleteTask(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusNoContent, gin.H{"message": "deleted task successfully"})
}

// AddTask handles POST /tasks
func AddTask(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		Status      string `json:"status"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	layout := "02-01-2006"
	parsedTime, err := time.Parse(layout, input.DueDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The time format is not 02-01-2006"})
		return
	}

	newTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     parsedTime,
		Status:      input.Status,
	}

	if newTask.Status != "Pending" && newTask.Status != "In Progress" && newTask.Status != "Completed" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid status value"})
		return
	}

	createdTask,err := data.AddTask(newTask)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to add task"})
		return
	}
	c.IndentedJSON(http.StatusCreated, createdTask)
}