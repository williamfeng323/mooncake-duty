package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//IDocumentBase Interface which each collection document (model) has to implement
type IDocumentBase interface {
	SetCollection(*mongo.Collection)
	Save() error
	Update(interface{}) (map[string]interface{}, error)
	Validate(...interface{}) (bool, []error)
}

// BaseModel is the base model that other models should embedded.
type BaseModel struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	Deleted    bool               `json:"-" bson:"deleted"`
}

// SetCollection will set the document's own collection instance
func (model *BaseModel) SetCollection(coll *mongo.Collection) {
	model.collection = coll
}

// Save is the default save function to persistence the document
func (model *BaseModel) Save() error {
	return nil
}

// Update is the default function to update the document
func (model *BaseModel) Update(interface{}) (map[string]interface{}, error) {
	return nil, nil
}

// Validate is the way you validate the document you save.
func (model *BaseModel) Validate(...interface{}) (bool, []error) {
	return true, nil
}
