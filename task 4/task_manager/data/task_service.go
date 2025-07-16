package data

import (
	"errors"
	"strconv"
	"time"

	"task_manager/models"
)

// In-memory task storage
var tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

// GetAllTasks returns all tasks
func GetAllTasks() []models.Task {
	return tasks
}

// GetTaskByID returns a task by ID
func GetTaskByID(id string) (*models.Task, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}

// UpdateTask updates a task by ID with new data
func UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i] = updatedTask
			tasks[i].ID = id 
			return &tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}

// DeleteTask deletes a task by ID
func DeleteTask(id string) error {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

// AddTask adds a new task and returns it
func AddTask(task models.Task) models.Task {
	newID := 1
	if len(tasks) > 0 {
		lastID, err := strconv.Atoi(tasks[len(tasks)-1].ID)
		if err == nil {
			newID = lastID + 1
		}
	}
	task.ID = strconv.Itoa(newID)
	tasks = append(tasks, task)
	return task
}
