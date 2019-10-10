package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//IDocumentBase Interface which each collection document (model) hast to implement
type IDocumentBase interface {
	Save() error
	Update(interface{}) (error, map[string]interface{})
	Validate(...interface{}) (bool, []error)
	DefaultValidate() (bool, []error)
}

// BaseModel the basic model that other models should embedded.
type BaseModel struct {
	document  IDocumentBase
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	Deleted   bool               `json:"-" bson:"deleted"`
}

//New return an model instance.
func (model *BaseModel) New() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return client, err
	}
	return client, nil
}
