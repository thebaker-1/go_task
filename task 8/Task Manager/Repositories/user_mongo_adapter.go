package Repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMongoCollectionAdapter struct {
	Coll *mongo.Collection
}

func (a *UserMongoCollectionAdapter) InsertOne(ctx context.Context, document interface{}, opts ...interface{}) (interface{}, error) {
	var mongoOpts []*options.InsertOneOptions
	for _, o := range opts {
		if opt, ok := o.(*options.InsertOneOptions); ok {
			mongoOpts = append(mongoOpts, opt)
		}
	}
	res, err := a.Coll.InsertOne(ctx, document, mongoOpts...)
	if err != nil {
		return nil, err
	}
	return &InsertOneResult{InsertedID: res.InsertedID}, nil
}

func (a *UserMongoCollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...interface{}) SingleResult {
	var mongoOpts []*options.FindOneOptions
	for _, o := range opts {
		if opt, ok := o.(*options.FindOneOptions); ok {
			mongoOpts = append(mongoOpts, opt)
		}
	}
	return a.Coll.FindOne(ctx, filter, mongoOpts...)
}
