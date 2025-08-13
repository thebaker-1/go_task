package Repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionAdapter struct {
	Coll *mongo.Collection
}

func (a *MongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...interface{}) (Cursor, error) {
	var mongoOpts []*options.FindOptions
	for _, o := range opts {
		if opt, ok := o.(*options.FindOptions); ok {
			mongoOpts = append(mongoOpts, opt)
		}
	}
	cursor, err := a.Coll.Find(ctx, filter, mongoOpts...)
	if err != nil {
		return nil, err
	}
	return cursor, nil // cursor implements Cursor
}

func (a *MongoCollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...interface{}) SingleResult {
	var mongoOpts []*options.FindOneOptions
	for _, o := range opts {
		if opt, ok := o.(*options.FindOneOptions); ok {
			mongoOpts = append(mongoOpts, opt)
		}
	}
	return a.Coll.FindOne(ctx, filter, mongoOpts...)
}

func (a *MongoCollectionAdapter) InsertOne(ctx context.Context, document interface{}, opts ...interface{}) (*InsertOneResult, error) {
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

func (a *MongoCollectionAdapter) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) SingleResult {
	var mongoOpts []*options.FindOneAndUpdateOptions
	for _, o := range opts {
		if opt, ok := o.(*options.FindOneAndUpdateOptions); ok {
			mongoOpts = append(mongoOpts, opt)
		}
	}
	return a.Coll.FindOneAndUpdate(ctx, filter, update, mongoOpts...)
}

func (a *MongoCollectionAdapter) DeleteOne(ctx context.Context, filter interface{}, opts ...interface{}) (DeleteResult, error) {
	var mongoOpts []*options.DeleteOptions
	for _, o := range opts {
		if opt, ok := o.(*options.DeleteOptions); ok {
			mongoOpts = append(mongoOpts, opt)
		}
	}
	res, err := a.Coll.DeleteOne(ctx, filter, mongoOpts...)
	if err != nil {
		return nil, err
	}
	return &MongoDeleteResultAdapter{res}, nil
}

type MongoDeleteResultAdapter struct {
	res *mongo.DeleteResult
}

func (m *MongoDeleteResultAdapter) DeletedCount() int64 {
	return m.res.DeletedCount
}
