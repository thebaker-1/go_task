package Domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// import "time"

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	DueDate     time.Time          `bson:"due_date"`
	Status      string             `bson:"status"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"` // Hashed password
	Email    string             `bson:"email"`
	Role     string             `bson:"role"` // e.g., "admin", "user"
}
