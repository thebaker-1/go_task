package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" binding:"required"`
	Password string             `json:"password" bson:"password" binding:"required"` // Store hashed password
	Email    string             `json:"email" bson:"email"`
	Role     string             `json:"role" bson:"role"` // e.g., "admin", "user"

}
