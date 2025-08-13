// Minimal, real test for Usecases.TaskUsecase using testify/mock
package tests

import (
	"context"
	"taskmanager/Domain"
	"taskmanager/Usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTaskRepository implements the TaskRepository interface for testing
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetAllTasks(ctx context.Context) ([]Domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) AddTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskUsecase_AddTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	inputTask := Domain.Task{Title: "Learn Go Usecases"}
	expectedTask := &Domain.Task{Title: "Learn Go Usecases"}

	mockRepo.On("AddTask", mock.Anything, inputTask).Return(expectedTask, nil)

	result, err := usecase.AddTask(context.Background(), inputTask)

	assert.NoError(t, err)
	assert.Equal(t, expectedTask, result)
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetAllTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	expectedTasks := []Domain.Task{{Title: "Task 1"}, {Title: "Task 2"}}
	mockRepo.On("GetAllTasks", mock.Anything).Return(expectedTasks, nil)

	result, err := usecase.GetAllTasks(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, result)
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTaskByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	fakeID := primitive.NewObjectID()
	expectedTask := &Domain.Task{ID: fakeID, Title: "Task by ID"}
	mockRepo.On("GetTaskByID", mock.Anything, fakeID).Return(expectedTask, nil)

	result, err := usecase.GetTaskByID(context.Background(), fakeID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, expectedTask, result)
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_UpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	fakeID := primitive.NewObjectID()
	inputTask := Domain.Task{Title: "Updated Task"}
	expectedTask := &Domain.Task{ID: fakeID, Title: "Updated Task"}
	// The usecase will set inputTask.ID = fakeID before calling UpdateTask
	mockRepo.On("UpdateTask", mock.Anything, Domain.Task{ID: fakeID, Title: "Updated Task"}).Return(expectedTask, nil)

	result, err := usecase.UpdateTask(context.Background(), fakeID.Hex(), inputTask)

	assert.NoError(t, err)
	assert.Equal(t, expectedTask, result)
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_DeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := Usecases.NewTaskUsecase(mockRepo)

	fakeID := primitive.NewObjectID()
	mockRepo.On("DeleteTask", mock.Anything, fakeID).Return(nil)

	err := usecase.DeleteTask(context.Background(), fakeID.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
