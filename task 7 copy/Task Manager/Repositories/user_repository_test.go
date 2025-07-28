package Repositories

import (
	"context"
	"errors"
	"testing"
	"taskmanager/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCollection simulates a MongoDB collection for testing
type MockCollection struct {
	data map[primitive.ObjectID]Domain.User
}

func NewMockCollection() *MockCollection {
	return &MockCollection{data: make(map[primitive.ObjectID]Domain.User)}
}

func (m *MockCollection) InsertOne(ctx context.Context, user Domain.User) error {
	for _, u := range m.data {
		if u.Username == user.Username {
			return errors.New("duplicate username")
		}
	}
	m.data[user.ID] = user
	return nil
}

func (m *MockCollection) FindOneByUsername(ctx context.Context, username string) (*Domain.User, error) {
	for _, u := range m.data {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockCollection) FindOneByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	user, exists := m.data[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// MongoUserRepository with mock collection
type TestMongoUserRepository struct {
	collection *MockCollection
}

func NewTestMongoUserRepository() *TestMongoUserRepository {
	return &TestMongoUserRepository{collection: NewMockCollection()}
}

func (r *TestMongoUserRepository) RegisterUser(ctx context.Context, user Domain.User) error {
	return r.collection.InsertOne(ctx, user)
}

func (r *TestMongoUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error) {
	user, err := r.collection.FindOneByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (r *TestMongoUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	return r.collection.FindOneByID(ctx, id)
}

func TestRegisterUser(t *testing.T) {
	repo := NewTestMongoUserRepository()

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	err := repo.RegisterUser(context.Background(), user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Duplicate username
	err = repo.RegisterUser(context.Background(), user)
	if err == nil {
		t.Errorf("expected error for duplicate username, got nil")
	}
}

func TestAuthenticateUser(t *testing.T) {
	repo := NewTestMongoUserRepository()

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	_ = repo.RegisterUser(context.Background(), user)

	// Correct credentials
	authUser, err := repo.AuthenticateUser(context.Background(), "testuser", "hashedpassword")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if authUser.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %v", authUser.Username)
	}

	// Incorrect password
	_, err = repo.AuthenticateUser(context.Background(), "testuser", "wrongpassword")
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}

	// Non-existent user
	_, err = repo.AuthenticateUser(context.Background(), "nouser", "password")
	if err == nil {
		t.Errorf("expected error for non-existent user, got nil")
	}
}

func TestGetUserByID(t *testing.T) {
	repo := NewTestMongoUserRepository()

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	_ = repo.RegisterUser(context.Background(), user)

	// Valid ID
	gotUser, err := repo.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if gotUser.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %v", gotUser.Username)
	}

	// Non-existent ID
	nonExistentID := primitive.NewObjectID()
	_, err = repo.GetUserByID(context.Background(), nonExistentID)
	if err == nil {
		t.Errorf("expected error for non-existent ID, got nil")
	}
}
