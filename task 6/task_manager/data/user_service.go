package data

import (
	"context"
	"errors"

	// "time"

	"task_mdb/models"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserService(client *mongo.Client, dbName string, collectionName string) *UserService {
	// Use background context without immediate cancel to avoid context canceled error
	ctx := context.Background()
	collection := client.Database(dbName).Collection(collectionName)

	// Create unique indexes on username and email fields
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		panic("Failed to create indexes: " + err.Error())
	}

	return &UserService{
		collection: collection,
		ctx:        ctx,
	}
}

func (s *UserService) RegisterUser(user *models.User) error {
	// Check if username exists
	count, err := s.collection.CountDocuments(s.ctx, bson.M{"username": user.Username})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	// Check if email exists
	count, err = s.collection.CountDocuments(s.ctx, bson.M{"email": user.Email})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	user.ID = primitive.NewObjectID()
	if user.Role == "" {
		user.Role = "user"
	}

	_, err = s.collection.InsertOne(s.ctx, user)
	if err != nil {
		// Check for duplicate key error
		if writeException, ok := err.(mongo.WriteException); ok {
			for _, we := range writeException.WriteErrors {
				if we.Code == 11000 {
					return errors.New("username or email already exists")
				}
			}
		}
		return err
	}
	return nil
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(s.ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(s.ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
