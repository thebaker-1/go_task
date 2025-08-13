package tests

import (
	"taskmanager/Domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test basic Task creation
func TestTaskCreation(t *testing.T) {
	task := Domain.Task{
		Title:   "Test Task",
		Status:  "pending",
		DueDate: time.Now().Add(24 * time.Hour),
	}
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "pending", task.Status)
	assert.NotZero(t, task.DueDate)
}

// Test basic User creation
func TestUserCreation(t *testing.T) {
	user := Domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "user", user.Role)
	assert.Contains(t, user.Email, "@")
}

// Table-driven test for valid/invalid Task statuses
func TestTaskStatusValidation(t *testing.T) {
	validStatuses := []string{"pending", "completed", "in-progress"}
	for _, status := range validStatuses {
		task := Domain.Task{Status: status}
		assert.Contains(t, validStatuses, task.Status)
	}

	invalidStatuses := []string{"unknown", "", "done"}
	for _, status := range invalidStatuses {
		task := Domain.Task{Status: status}
		assert.NotContains(t, validStatuses, task.Status)
	}
}

// Test Task due date in the past
func TestTaskDueDateInPast(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	task := Domain.Task{DueDate: pastDate}
	assert.True(t, task.DueDate.Before(time.Now()))
}

// Test Task due date in the future
func TestTaskDueDateInFuture(t *testing.T) {
	futureDate := time.Now().Add(48 * time.Hour)
	task := Domain.Task{DueDate: futureDate}
	assert.True(t, task.DueDate.After(time.Now()))
}

// Example: Test User email validation (if you add a method to Domain.User)
func TestUserEmailValidation(t *testing.T) {
	validUser := Domain.User{Email: "valid@example.com"}
	invalidUser := Domain.User{Email: "invalid-email"}
	assert.True(t, validUser.IsValidEmail())
	assert.False(t, invalidUser.IsValidEmail())
}

func TestTaskIsOverdueV1(t *testing.T) {
	overdueTask := Domain.Task{DueDate: time.Now().Add(-1 * time.Hour)}
	futureTask := Domain.Task{DueDate: time.Now().Add(1 * time.Hour)}
	assert.True(t, overdueTask.IsOverdue())
	assert.False(t, futureTask.IsOverdue())
}

func TestTaskIsOverdueV2(t *testing.T) {
	cases := []struct {
		name     string
		dueDate  time.Time
		expected bool
	}{
		{"Overdue task", time.Now().Add(-1 * time.Hour), true},
		{"Future task", time.Now().Add(1 * time.Hour), false},
		{"Due now", time.Now(), false}, // assuming IsOverdue returns false if due now
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			task := Domain.Task{DueDate: tc.dueDate}
			assert.Equal(t, tc.expected, task.IsOverdue())
		})
	}
}
