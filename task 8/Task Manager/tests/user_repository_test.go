package tests

import (
	"context"
	// "errors"
	"taskmanager/Domain"
	"taskmanager/Repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Tests ---
// --- Unique mocks for user repository tests to avoid naming conflicts ---
type UserMockCollection struct{ mock.Mock }

func (m *UserMockCollection) InsertOne(ctx context.Context, doc interface{}, opts ...interface{}) (interface{}, error) {
	args := m.Called(ctx, doc)
	return args.Get(0), args.Error(1)
}

func (m *UserMockCollection) FindOne(ctx context.Context, filter interface{}, opts ...interface{}) Repositories.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(Repositories.SingleResult)
}

type UserMockSingleResult struct {
	user Domain.User
	err  error
}

func (m *UserMockSingleResult) Decode(val interface{}) error {
	if m.err != nil {
		return m.err
	}
	*(val.(*Domain.User)) = m.user
	return nil
}

// 	return nil
// }

func TestRegisterUser(t *testing.T) {
	mockColl := new(UserMockCollection)
	repo := Repositories.NewMongoUserRepository(mockColl)
	user := Domain.User{Username: "testuser"}
	mockColl.On("InsertOne", mock.Anything, user).Return(nil, nil)
	err := repo.RegisterUser(context.Background(), user)
	assert.NoError(t, err)
	mockColl.AssertExpectations(t)
}

func TestAuthenticateUser_Found(t *testing.T) {
	mockColl := new(UserMockCollection)
	repo := Repositories.NewMongoUserRepository(mockColl)
	user := Domain.User{Username: "testuser", Password: "pass"}
	mockResult := &UserMockSingleResult{user: user}
	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)
	result, err := repo.AuthenticateUser(context.Background(), "testuser", "pass")
	assert.NoError(t, err)
	assert.Equal(t, &user, result)
	mockColl.AssertExpectations(t)
}

func TestAuthenticateUser_NotFound(t *testing.T) {
	mockColl := new(UserMockCollection)
	repo := Repositories.NewMongoUserRepository(mockColl)
	mockResult := &UserMockSingleResult{err: assert.AnError}
	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)
	result, err := repo.AuthenticateUser(context.Background(), "testuser", "pass")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockColl.AssertExpectations(t)
}

func TestGetUserByID_Found(t *testing.T) {
	mockColl := new(UserMockCollection)
	repo := Repositories.NewMongoUserRepository(mockColl)
	fakeID := struct{ primitive.ObjectID }{primitive.NewObjectID()}.ObjectID
	user := Domain.User{ID: fakeID, Username: "testuser"}
	mockResult := &UserMockSingleResult{user: user}
	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)
	result, err := repo.GetUserByID(context.Background(), fakeID)
	assert.NoError(t, err)
	assert.Equal(t, &user, result)
	mockColl.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockColl := new(UserMockCollection)
	repo := Repositories.NewMongoUserRepository(mockColl)
	fakeID := struct{ primitive.ObjectID }{primitive.NewObjectID()}.ObjectID
	mockResult := &UserMockSingleResult{err: assert.AnError}
	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)
	result, err := repo.GetUserByID(context.Background(), fakeID)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockColl.AssertExpectations(t)
}
