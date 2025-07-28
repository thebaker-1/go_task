package controllers

import (
	"context"
	"net/http"
	"time"

	"taskmanager/Domain"
	"taskmanager/Usecases"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define your JWT secret key. In a real app, this should be an environment variable.

var jwtSecret []byte

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set in .env file")
	}
	jwtSecret = []byte(secret)
}

type Controller struct {
	UserUsecase Usecases.UserUsecase
	TaskUsecase Usecases.TaskUsecase
}

type TaskDTO struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	DueDate     string `json:"due_date" bson:"due_date"`
	Status      string `json:"status" bson:"status"`
}

type UserDTO struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
	Email    string `json:"email" bson:"email"`
	Role     string `json:"role" bson:"role"`
}

func NewController(userUsecase Usecases.UserUsecase, taskUsecase Usecases.TaskUsecase) *Controller {
	return &Controller{
		UserUsecase: userUsecase,
		TaskUsecase: taskUsecase,
	}
}

func toTaskDTO(task Domain.Task) TaskDTO {
	return TaskDTO{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate.Format("02-01-2006"),
		Status:      task.Status,
	}
}

func toTaskDomain(dto TaskDTO) (Domain.Task, error) {
	dueDate, err := time.Parse("02-01-2006", dto.DueDate)
	if err != nil {
		return Domain.Task{}, err
	}
	var id primitive.ObjectID
	if dto.ID != "" {
		id, err = primitive.ObjectIDFromHex(dto.ID)
		if err != nil {
			return Domain.Task{}, err
		}
	} else {
		id = primitive.NilObjectID
	}
	return Domain.Task{
		ID:          id,
		Title:       dto.Title,
		Description: dto.Description,
		DueDate:     dueDate,
		Status:      dto.Status,
	}, nil
}

// func toUserDTO(user Domain.User) UserDTO {
// 	// This function is currently unused but kept for future use or consistency
// 	return UserDTO{
// 		ID:       user.ID.Hex(),
// 		Username: user.Username,
// 		Password: user.Password,
// 		Email:    user.Email,
// 		Role:     user.Role,
// 	}
// }

func toUserDomain(dto UserDTO) Domain.User {
	var id primitive.ObjectID
	// Since UserDTO.ID is string, convert to ObjectID if not empty
	if dto.ID != "" {
		var err error
		id, err = primitive.ObjectIDFromHex(dto.ID)
		if err != nil {
			id = primitive.NilObjectID
		}
	} else {
		id = primitive.NilObjectID
	}
	return Domain.User{
		ID:       id,
		Username: dto.Username,
		Password: dto.Password,
		Email:    dto.Email,
		Role:     dto.Role,
	}
}

func (c *Controller) RegisterUser(ctx *gin.Context) {
	var newUserDTO UserDTO
	if err := ctx.ShouldBindJSON(&newUserDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := toUserDomain(newUserDTO)

	if err := c.UserUsecase.RegisterUser(context.Background(), newUser); err != nil {
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

	user, err := c.UserUsecase.AuthenticateUser(context.Background(), loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
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
func (c *Controller) GetTasks(ctx *gin.Context) {
	tasks, err := c.TaskUsecase.GetAllTasks(context.Background())
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "1Failed to retrieve tasks", "error": err.Error()})
		return
	}

	if len(tasks) == 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "There are no tasks yet"})
		return
	}

	var taskDTOs []TaskDTO
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, toTaskDTO(task))
	}

	ctx.IndentedJSON(http.StatusOK, taskDTOs)
}

// GetTask handles GET /tasks/:id
func (c *Controller) GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.TaskUsecase.GetTaskByID(context.Background(), id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, toTaskDTO(*task))
}

// UpdateTask handles PUT /tasks/:id
func (c *Controller) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var input TaskDTO

	if err := ctx.BindJSON(&input); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	task, err := toTaskDomain(input)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid due date format"})
		return
	}
	// Convert string id to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}
	task.ID = objID

	updatedTask, err := c.TaskUsecase.UpdateTask(context.Background(), id, task)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, toTaskDTO(*updatedTask))
}

// DeleteTask handles DELETE /tasks/:id
func (c *Controller) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.TaskUsecase.DeleteTask(context.Background(), id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	// Return 200 OK with message instead of 204 No Content with body
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "deleted task successfully"})
}

// AddTask handles POST /tasks
func (c *Controller) AddTask(ctx *gin.Context) {
	var input TaskDTO

	if err := ctx.BindJSON(&input); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	task, err := toTaskDomain(input)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid due date format"})
		return
	}

	if task.Status != "Pending" && task.Status != "In Progress" && task.Status != "Completed" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid status value"})
		return
	}

	createdTask, err := c.TaskUsecase.AddTask(context.Background(), task)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to add task"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, toTaskDTO(*createdTask))
}
