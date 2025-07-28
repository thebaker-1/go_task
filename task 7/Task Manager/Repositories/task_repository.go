package Repositories

import (
	"context"
	"errors"
	"fmt"
	"taskmanager/Domain"
	// import "go.mongodb.org/mongo-driver/bson/primitive"
	
	"go.mongodb.org/mongo-driver/mongo/options"
	// MongoTaskRepository implements TaskRepository using MongoDB
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskEntity is the persistence model for Task with BSON tags and ObjectID
type TaskEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	DueDate     primitive.DateTime `bson:"due_date"`
	Status      string             `bson:"status"`
}


// ToDomain converts TaskEntity to Domain.Task
func (te *TaskEntity) ToDomain() Domain.Task {
	return Domain.Task{
		ID:          te.ID,
		Title:       te.Title,
		Description: te.Description,
		DueDate:     te.DueDate.Time(),
		Status:      te.Status,
	}
}

// FromDomain converts Domain.Task to TaskEntity
func FromDomain(task Domain.Task) TaskEntity {
	return TaskEntity{
		ID:          task.ID, // Use domain model ID directly
		Title:       task.Title,
		Description: task.Description,
		DueDate:     primitive.NewDateTimeFromTime(task.DueDate),
		Status:      task.Status,
	}
}

type TaskRepository interface {
	GetAllTasks(ctx context.Context) ([]Domain.Task, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Domain.Task, error)
	AddTask(ctx context.Context, task Domain.Task) (*Domain.Task, error)
	UpdateTask(ctx context.Context, task Domain.Task) (*Domain.Task, error)
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
}

type MongoTaskRepository struct {
	collection *mongo.Collection
}

// Cursor interface abstracts mongo.Cursor methods used
type Cursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
	Err() error
}

// SingleResult interface abstracts mongo.SingleResult methods used
type SingleResult interface {
	Decode(interface{}) error
}

// InsertOneResult interface abstracts mongo.InsertOneResult
type InsertOneResult struct {
	InsertedID interface{}
}

// DeleteResult interface abstracts mongo.DeleteResult
type DeleteResult interface {
	DeletedCount() int64
}

func NewMongoTaskRepository(collection *mongo.Collection) *MongoTaskRepository {
	return &MongoTaskRepository{collection: collection}
}

func (r *MongoTaskRepository) GetAllTasks(ctx context.Context) ([]Domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var taskEntities []TaskEntity
	for cursor.Next(ctx) {
		var taskEntity TaskEntity
		if err := cursor.Decode(&taskEntity); err != nil {
			return nil, fmt.Errorf("error decoding task entity: %w", err)
		}
		taskEntities = append(taskEntities, taskEntity)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	var tasks []Domain.Task
	for _, te := range taskEntities {
		tasks = append(tasks, te.ToDomain())
	}

	return tasks, nil
}

func (r *MongoTaskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Domain.Task, error) {
	var taskEntity TaskEntity
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&taskEntity)
	if err != nil {
		return nil, errors.New("task not found")
	}
	task := taskEntity.ToDomain()
	return &task, nil
}

func (r *MongoTaskRepository) AddTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	taskEntity := FromDomain(task)
	taskEntity.ID = primitive.NilObjectID
	res, err := r.collection.InsertOne(ctx, taskEntity)
	if err != nil {
		return nil, err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}
	taskEntity.ID = id
	domainTask := taskEntity.ToDomain()
	return &domainTask, nil
}


func (r *MongoTaskRepository) UpdateTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	taskEntity := FromDomain(task)
	filter := bson.M{"_id": taskEntity.ID}
	update := bson.M{"$set": taskEntity}
	var updatedTaskEntity TaskEntity
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedTaskEntity)
	if err != nil {
		return nil, errors.New("task not found")
	}
	updatedTask := updatedTaskEntity.ToDomain()
	return &updatedTask, nil
}

func (r *MongoTaskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
