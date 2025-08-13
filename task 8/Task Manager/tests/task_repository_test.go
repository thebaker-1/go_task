package tests

import (
	"context"
	"taskmanager/Domain"
	"taskmanager/Repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// // MockCollection is a simple mock for the Collection interface
type MockCollection struct{ mock.Mock }

// Add a stub InsertOne to satisfy the interface
func (m *MockCollection) InsertOne(ctx context.Context, doc interface{}, opts ...interface{}) (*Repositories.InsertOneResult, error) {
	args := m.Called(ctx, doc)
	return args.Get(0).(*Repositories.InsertOneResult), args.Error(1)
}

// --- Minimal mocks for SingleResult and DeleteResult ---
type MockSingleResult struct {
	entity Repositories.TaskEntity
	err    error
}

// Add a stub DeleteOne to satisfy the interface
func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...interface{}) (Repositories.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(Repositories.DeleteResult), args.Error(1)
}

// Add a stub Find to satisfy the interface
func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...interface{}) (Repositories.Cursor, error) {
	return nil, nil
}

// Add a stub FindOne to satisfy the interface
func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...interface{}) Repositories.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(Repositories.SingleResult)
}

// Add a stub FindOneAndUpdate to satisfy the interface
func (m *MockCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) Repositories.SingleResult {
	return nil
}

func (m *MockSingleResult) Decode(val interface{}) error {
	if m.err != nil {
		return m.err
	}
	*(val.(*Repositories.TaskEntity)) = m.entity
	return nil
}

type MockDeleteResult struct {
	deleted int64
}

func (m *MockDeleteResult) DeletedCount() int64 {
	return m.deleted
}

func TestMongoTaskRepository_AddTask(t *testing.T) {
	mockColl := new(MockCollection)
	repo := Repositories.NewMongoTaskRepository(mockColl)

	// Prepare input and expected output
	domainTask := Domain.Task{
		Title: "Minimal Real Test",
	}
	insertedID := primitive.NewObjectID()
	// The repository expects InsertOne to return this result
	mockColl.On("InsertOne", mock.Anything, mock.Anything).Return(&Repositories.InsertOneResult{InsertedID: insertedID}, nil)

	// Call the real repository method
	result, err := repo.AddTask(context.Background(), domainTask)

	// Check the result
	assert.NoError(t, err)
	assert.Equal(t, "Minimal Real Test", result.Title)
	assert.Equal(t, insertedID, result.ID)

	mockColl.AssertExpectations(t)
}

func TestMongoTaskRepository_GetTaskByID_Found(t *testing.T) {
	mockColl := new(MockCollection)
	repo := Repositories.NewMongoTaskRepository(mockColl)
	fakeID := primitive.NewObjectID()
	entity := Repositories.TaskEntity{ID: fakeID, Title: "Found Task"}
	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(&MockSingleResult{entity: entity, err: nil})

	result, err := repo.GetTaskByID(context.Background(), fakeID)
	assert.NoError(t, err)
	assert.Equal(t, "Found Task", result.Title)
	assert.Equal(t, fakeID, result.ID)
	mockColl.AssertExpectations(t)
}

func TestMongoTaskRepository_GetTaskByID_NotFound(t *testing.T) {
	mockColl := new(MockCollection)
	repo := Repositories.NewMongoTaskRepository(mockColl)
	fakeID := primitive.NewObjectID()
	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(&MockSingleResult{err: assert.AnError})

	result, err := repo.GetTaskByID(context.Background(), fakeID)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockColl.AssertExpectations(t)
}

func TestMongoTaskRepository_DeleteTask_Found(t *testing.T) {
	mockColl := new(MockCollection)
	repo := Repositories.NewMongoTaskRepository(mockColl)
	fakeID := primitive.NewObjectID()
	mockColl.On("DeleteOne", mock.Anything, mock.Anything).Return(&MockDeleteResult{deleted: 1}, nil)

	err := repo.DeleteTask(context.Background(), fakeID)
	assert.NoError(t, err)
	mockColl.AssertExpectations(t)
}

func TestMongoTaskRepository_DeleteTask_NotFound(t *testing.T) {
	mockColl := new(MockCollection)
	repo := Repositories.NewMongoTaskRepository(mockColl)
	fakeID := primitive.NewObjectID()
	mockColl.On("DeleteOne", mock.Anything, mock.Anything).Return(&MockDeleteResult{deleted: 0}, nil)

	err := repo.DeleteTask(context.Background(), fakeID)
	assert.Error(t, err)
	mockColl.AssertExpectations(t)
}
