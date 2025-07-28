package Usecases

import (
	"context"
	"errors"
	"taskmanager/Domain"
	"taskmanager/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskUsecase defines the use case interface for task operations
type TaskUsecase interface {
	GetAllTasks(ctx context.Context) ([]Domain.Task, error)
	GetTaskByID(ctx context.Context, id string) (*Domain.Task, error)
	AddTask(ctx context.Context, task Domain.Task) (*Domain.Task, error)
	UpdateTask(ctx context.Context, id string, task Domain.Task) (*Domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

// taskUsecase implements TaskUsecase interface
type taskUsecase struct {
	taskRepo Repositories.TaskRepository
}

// NewTaskUsecase creates a new TaskUsecase
func NewTaskUsecase(taskRepo Repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{taskRepo: taskRepo}
}

func (u *taskUsecase) GetAllTasks(ctx context.Context) ([]Domain.Task, error) {
	tasks, err := u.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *taskUsecase) GetTaskByID(ctx context.Context, id string) (*Domain.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}
	return u.taskRepo.GetTaskByID(ctx, objID)
}

func (u *taskUsecase) AddTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	createdTask, err := u.taskRepo.AddTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return createdTask, nil
}

func (u *taskUsecase) UpdateTask(ctx context.Context, id string, task Domain.Task) (*Domain.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}
	task.ID = objID
	return u.taskRepo.UpdateTask(ctx, task)
}

func (u *taskUsecase) DeleteTask(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}
	return u.taskRepo.DeleteTask(ctx, objID)
}
