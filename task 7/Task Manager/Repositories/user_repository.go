package Repositories

import (
	"context"
	"errors"
	"taskmanager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	RegisterUser(ctx context.Context, user Domain.User) error
	AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error)
}

// MongoUserRepository implements UserRepository using MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoUserRepository
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) RegisterUser(ctx context.Context, user Domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error) {
	var user Domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *MongoUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	var user Domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
