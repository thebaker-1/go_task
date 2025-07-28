package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"taskmanager/Domain"
	"taskmanager/Usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserUsecase implements Usecases.UserUsecase for testing
type MockUserUsecase struct {
	users map[string]Domain.User
}

func NewMockUserUsecase() *MockUserUsecase {
	return &MockUserUsecase{users: make(map[string]Domain.User)}
}

func (m *MockUserUsecase) RegisterUser(ctx context.Context, user Domain.User) error {
	if _, exists := m.users[user.Username]; exists {
		return &UserExistsError{}
	}
	m.users[user.Username] = user
	return nil
}

func (m *MockUserUsecase) AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, &UserNotFoundError{}
	}
	if user.Password != password {
		return nil, &InvalidPasswordError{}
	}
	return &user, nil
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, id string) (*Domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	for _, user := range m.users {
		if user.ID == objID {
			return &user, nil
		}
	}
	return nil, &UserNotFoundError{}
}

type UserExistsError struct{}

func (e *UserExistsError) Error() string {
	return "username already exists"
}

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "user not found"
}

type InvalidPasswordError struct{}

func (e *InvalidPasswordError) Error() string {
	return "invalid password"
}

func setupRouter(userUsecase Usecases.UserUsecase) *gin.Engine {
	r := gin.Default()
	ctrl := NewController(userUsecase, nil)

	r.POST("/register", ctrl.RegisterUser)
	r.POST("/login", ctrl.LoginUser)
	return r
}

func TestRegisterUser(t *testing.T) {
	userUsecase := NewMockUserUsecase()
	router := setupRouter(userUsecase)

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Try to register the same user again
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409 for duplicate user, got %d", w.Code)
	}
}

func TestLoginUser(t *testing.T) {
	userUsecase := NewMockUserUsecase()
	router := setupRouter(userUsecase)

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	_ = userUsecase.RegisterUser(context.Background(), user)

	loginPayload := map[string]string{
		"username": "testuser",
		"password": "hashedpassword",
	}
	jsonValue, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Incorrect password
	loginPayload["password"] = "wrongpassword"
	jsonValue, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for wrong password, got %d", w.Code)
	}
}

func TestGetUserByID(t *testing.T) {
	userUsecase := NewMockUserUsecase()
	router := setupRouter(userUsecase)

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "testuser@example.com",
		Role:     "user",
	}

	_ = userUsecase.RegisterUser(context.Background(), user)

	req, _ := http.NewRequest("GET", "/user/"+user.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Invalid ID
	req, _ = http.NewRequest("GET", "/user/invalidid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid ID, got %d", w.Code)
	}

	// Non-existent ID
	req, _ = http.NewRequest("GET", "/user/507f1f77bcf86cd799439011", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent user, got %d", w.Code)
	}
}
