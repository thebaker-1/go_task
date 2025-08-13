package Domain

// import "time"

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

// IsOverdue checks if the task's due date is in the past.
func (t *Task) IsOverdue() bool {
	return time.Now().After(t.DueDate)
}
