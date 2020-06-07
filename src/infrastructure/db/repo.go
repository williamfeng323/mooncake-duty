package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Validator helps validating the files
type Validator interface {
	// Verify returns empty error list if validate successfully.
	Verify(interface{}) []error
}

// Repository describe how to interact with database
type Repository interface {
	GetName() string
	SetCollection(*mongo.Collection)
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}
