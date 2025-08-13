package tests

import (
	"context"
	"errors"
	"taskmanager/Domain"
	"taskmanager/Usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository implements the UserRepository interface for usecase tests
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) RegisterUser(ctx context.Context, user Domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*Domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func TestRegisterUser_Usecase(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := Usecases.NewUserUsecase(mockRepo)
	user := Domain.User{Username: "testuser"}
	mockRepo.On("RegisterUser", mock.Anything, user).Return(nil)
	err := usecase.RegisterUser(context.Background(), user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_Usecase_Found(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := Usecases.NewUserUsecase(mockRepo)
	user := &Domain.User{Username: "testuser", Password: "pass"}
	mockRepo.On("AuthenticateUser", mock.Anything, "testuser", "pass").Return(user, nil)
	result, err := usecase.AuthenticateUser(context.Background(), "testuser", "pass")
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_Usecase_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := Usecases.NewUserUsecase(mockRepo)
	mockRepo.On("AuthenticateUser", mock.Anything, "testuser", "pass").Return(nil, errors.New("not found"))
	result, err := usecase.AuthenticateUser(context.Background(), "testuser", "pass")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Usecase_Found(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := Usecases.NewUserUsecase(mockRepo)
	fakeID := primitive.NewObjectID()
	user := &Domain.User{Username: "testuser"}
	mockRepo.On("GetUserByID", mock.Anything, fakeID).Return(user, nil)
	result, err := usecase.GetUserByID(context.Background(), fakeID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Usecase_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := Usecases.NewUserUsecase(mockRepo)
	fakeID := primitive.NewObjectID()
	mockRepo.On("GetUserByID", mock.Anything, fakeID).Return(nil, errors.New("not found"))
	result, err := usecase.GetUserByID(context.Background(), fakeID.Hex())
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
