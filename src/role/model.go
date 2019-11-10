package role

import (
	"context"
	"time"
	"williamfeng323/mooncake-duty/src/dao"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Permission the predefined permissions strings
type Permission string

const (
	// Create predefined permission mode, create permission
	Create Permission = "create"
	// Read predefined permission mode, read permission
	Read Permission = "read"
	// Update predefined permission mode, update permission
	Update Permission = "update"
	// Delete predefined permission mode, deleted permission
	Delete Permission = "delete"
)

//Role the role struct type for account.
type Role struct {
	dao.BaseModel `json:",inline" bson:",inline"`
	Name          string     `json:"name" bson:"name" required:"true"`
	Auth          Permission `json:"auth" bson:"auth" required:"true"`
}

// InsertRole create a new document in mongoDB with
// initial createdAt and _id
func (rol *Role) InsertRole() (*mongo.InsertOneResult, error) {
	validationErrors := rol.DefaultValidator()
	if validationErrors != nil {
		return nil, validationErrors[0]
	}
	rol.CreatedAt = time.Now()
	rol.ID = primitive.NewObjectID()
	bRole, err := bson.Marshal(rol)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return rol.GetCollection().InsertOne(ctx, bRole)
}
