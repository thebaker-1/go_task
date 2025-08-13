package Domain

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID
	Username string
	Password string // Hashed password
	Email    string
	Role     string // e.g., "admin", "user"
}

// IsValidEmail checks if the user's email is valid.
func (u *User) IsValidEmail() bool {
	// Simple regex for demonstration; you can use a more robust one.
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(u.Email)
}
