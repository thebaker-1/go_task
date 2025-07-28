package Usecases

import (
	"context"
	"errors"
	"testing"
	"taskmanager/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[string]Domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]Domain.User)}
}

func (m *MockUserRepository) RegisterUser(ctx context.Context, user Domain.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	if user.Password != password {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func TestRegisterUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	err := usecase.RegisterUser(context.Background(), user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Try to register the same user again
	err = usecase.RegisterUser(context.Background(), user)
	if err == nil {
		t.Errorf("expected error for duplicate user, got nil")
	}
}

func TestAuthenticateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	_ = usecase.RegisterUser(context.Background(), user)

	// Correct credentials
	authUser, err := usecase.AuthenticateUser(context.Background(), "testuser", "hashedpassword")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if authUser.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %v", authUser.Username)
	}

	// Incorrect password
	_, err = usecase.AuthenticateUser(context.Background(), "testuser", "wrongpassword")
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}

	// Non-existent user
	_, err = usecase.AuthenticateUser(context.Background(), "nouser", "password")
	if err == nil {
		t.Errorf("expected error for non-existent user, got nil")
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	_ = usecase.RegisterUser(context.Background(), user)

	// Valid ID
	gotUser, err := usecase.GetUserByID(context.Background(), user.ID.Hex())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if gotUser.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %v", gotUser.Username)
	}

	// Invalid ID
	_, err = usecase.GetUserByID(context.Background(), "invalidid")
	if err == nil {
		t.Errorf("expected error for invalid ID, got nil")
	}

	// Non-existent ID
	nonExistentID := primitive.NewObjectID()
	_, err = usecase.GetUserByID(context.Background(), nonExistentID.Hex())
	if err == nil {
		t.Errorf("expected error for non-existent ID, got nil")
	}
}
