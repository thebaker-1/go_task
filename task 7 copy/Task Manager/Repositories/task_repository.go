package Repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"taskmanager/Domain"

	// MongoTaskRepository implements TaskRepository using MongoDB
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskRepository defines the interface for task data access
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
	log.Println("GetAllTasks: Attempting to retrieve all tasks from MongoDB.")

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("GetAllTasks: Error finding documents: %v", err)
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("GetAllTasks: Error closing cursor: %v", cerr)
		}
		log.Println("GetAllTasks: Cursor closed.")
	}()

	var tasks []Domain.Task
	taskCount := 0
	for cursor.Next(ctx) {
		var task Domain.Task
		if err := cursor.Decode(&task); err != nil {
			// This is where your error "error decoding key id: an ObjectID string must be exactly 12 bytes long (got 0)"
			// is likely occurring. Log the problematic data if possible.
			log.Printf("GetAllTasks: Error decoding task. Potentially malformed document. Error: %v", err)
			// If you want to log the raw BSON document that caused the error (requires more effort):
			// var raw bson.Raw
			// if err := cursor.Decode(&raw); err == nil {
			//     log.Printf("GetAllTasks: Raw BSON document that caused decode error: %s", raw.String())
			// } else {
			//     log.Printf("GetAllTasks: Could not decode to raw BSON either: %v", err)
			// }

			// Depending on your requirements, you might choose to skip this malformed document
			// and continue, or return an error immediately. For now, we'll return an error.
			return nil, fmt.Errorf("error decoding task: %w", err)
		}
		tasks = append(tasks, task)
		taskCount++
		log.Printf("GetAllTasks: Successfully decoded task %d. Task ID: %s", taskCount, task.ID.Hex()) // Assuming Task.ID is a primitive.ObjectID
	}

	if err := cursor.Err(); err != nil {
		log.Printf("GetAllTasks: Cursor iteration error: %v", err)
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	log.Printf("GetAllTasks: Successfully retrieved %d tasks.", len(tasks))
	return tasks, nil
}

func (r *MongoTaskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*Domain.Task, error) {
	var task Domain.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (r *MongoTaskRepository) AddTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	task.ID = primitive.NilObjectID
	// Debug log before insertion
	// log.Printf("Inserting task: %+v\n", task)
	res, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		// log.Printf("InsertOne error: %v\n", err)
		return nil, err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		// log.Printf("Failed to convert inserted ID to ObjectID")
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}
	task.ID = id
	// log.Printf("Inserted task ID: %v\n", id.Hex())
	return &task, nil
}

func (r *MongoTaskRepository) UpdateTask(ctx context.Context, task Domain.Task) (*Domain.Task, error) {
	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": task}
	var updatedTask Domain.Task
	err := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedTask)
	if err != nil {
		return nil, errors.New("task not found")
	}
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
