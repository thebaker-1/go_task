package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
)

func InitMongo(uri string, dbName string, collectionName string) error {
	var err error
	Ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	Client, err = mongo.Connect(Ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	err = Client.Ping(Ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func DisconnectMongo() error {
	Cancel()
	return Client.Disconnect(Ctx)
}