package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

// GetTasks handles GET /tasks
func GetTasks(c *gin.Context) {
	tasks := data.GetAllTasks()
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

	createdTask := data.AddTask(newTask)
	c.IndentedJSON(http.StatusCreated, createdTask)
}
