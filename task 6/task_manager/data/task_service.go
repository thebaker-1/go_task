package data

import (
	"context"
	"errors"

	// "strconv"
	// "time"

	"task_mdb/models"

	// "github.com/go-playground/locales/id"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
    collection *mongo.Collection
)

func InitTaskService(client *mongo.Client, dbName string, collectionName string) error {
    if client == nil {
        return errors.New("mongo client is nil")
    }
    collection = client.Database(dbName).Collection(collectionName)
    return nil
}


// In-memory task storage
// var tasks = []models.Task{
// 	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
// 	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
// 	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
// }

// GetAllTasks returns all tasks
// func GetAllTasks() []models.Task {
// 	return tasks
// }

func GetAllTasks() ([]models.Task, error) {
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task
	for cursor.Next(context.TODO()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
// GetTaskByID returns a task by ID
// func GetTaskByID(id string) (*models.Task, error) {
// 	for i := range tasks {
// 		if tasks[i].ID == id {
// 			return &tasks[i], nil
// 		}
// 	}
// 	return nil, errors.New("task not found")
// }

func GetTaskByID(id string) (*models.Task, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid task ID")
    }
    var task models.Task
    err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&task)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("task not found")
        }
        return nil, err
    }
    return &task, nil
}


// UpdateTask updates a task by ID with new data
// func UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
// 	for i := range tasks {
// 		if tasks[i].ID == id {
// 			tasks[i] = updatedTask
// 			tasks[i].ID = id 
// 			return &tasks[i], nil
// 		}
// 	}
// 	return nil, errors.New("task not found")
// }
func UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid task ID")
    }
    updatedTask.ID = objID
    filter := bson.M{"_id": objID}
    update := bson.M{"$set": updatedTask}
    opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
    var result models.Task
    err = collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("task not found")
        }
        return nil, err
    }
    return &result, nil
}

// DeleteTask deletes a task by ID
// func DeleteTask(id string) error {
// 	for i := range tasks {
// 		if tasks[i].ID == id {
// 			tasks = append(tasks[:i], tasks[i+1:]...)
// 			return nil
// 		}
// 	}
// 	return errors.New("task not found")
// }
func DeleteTask(id string) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid task ID")
    }
    res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return errors.New("task not found")
    }
    return nil
}


// AddTask adds a new task and returns it
// func AddTask(task models.Task) models.Task {
// 	newID := 1
// 	if len(tasks) > 0 {
// 		lastID, err := strconv.Atoi(tasks[len(tasks)-1].ID)
// 		if err == nil {
// 			newID = lastID + 1
// 		}
// 	}
// 	task.ID = strconv.Itoa(newID)
// 	tasks = append(tasks, task)
// 	return task
// }
func AddTask(task models.Task) (*models.Task, error) {
    task.ID = primitive.NilObjectID // clear ID to let MongoDB generate ObjectID
    res, err := collection.InsertOne(context.TODO(), task)
    if err != nil {
        return nil, err
    }
    objID, ok := res.InsertedID.(primitive.ObjectID)
    if !ok {
        return nil, errors.New("failed to convert inserted ID to ObjectID")
    }
    task.ID = objID
    return &task, nil
}