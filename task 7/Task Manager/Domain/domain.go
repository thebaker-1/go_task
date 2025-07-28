package Domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

type User struct {
	ID       primitive.ObjectID
	Username string
	Password string // Hashed password
	Email    string
	Role     string // e.g., "admin", "user"
}
