package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel is the base model that other models should embedded.
// Please use tag `json:",inline" bson:",inline"` to make the exported
// fields inline.
type BaseModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt *time.Time         `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time         `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Deleted   bool               `json:"-" bson:"deleted,omitempty"`
}
