package project

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	dao "williamfeng323/mooncake-duty/src/infrastructure/db"
	"williamfeng323/mooncake-duty/src/domains/account"
)

// Permission the predefined permissions strings
type Permission string

// Tier is the support level of a team member.
type Tier string

const (
	// Create predefined permission mode, create permission
	Create Permission = "create"
	// Read predefined permission mode, read permission
	Read Permission = "read"
	// Update predefined permission mode, update permission
	Update Permission = "update"
	// Delete predefined permission mode, deleted permission
	Delete Permission = "delete"
	// T1 predefined tier level
	T1 Tier = "T1"
	// T2 predefined tier level
	T2 Tier = "T2"
	// T3 predefined tier level
	T3 Tier = "T3"
)

//Role the role struct type for account.
type Role struct {
	Name        string       `json:"name" bson:"name" required:"true"`
	Authorities []Permission `json:"authorities" bson:"authorities" required:"true"`
}

// Member is the struct to contain member information
type Member struct {
	AccountID primitive.ObjectID `json:"accountId" bson:"accountId"`
	Role      Role               `json:"role" bson:"role"`
	Tier      Tier               `json:"tier" bson:"tier"`
}

// Project is the struct to contain project information.
type Project struct {
	dao.BaseModel 									`json:",inline" bson:",inline"`
	Name          string   					`json:"name" bson:"name" required:"true"`
	Description   string   					`json:"description" bson:"description" required:"true"`
	Members       []account.Account `json:"members" bson:"members"`
}

// CreateProject create a new project document in mongoDB with
// initial createdAt and _id
func (project *Project) CreateProject() (*mongo.InsertOneResult, error) {
	validationErrors := project.DefaultValidator()
	if validationErrors != nil {
		return nil, validationErrors[0]
	}
	project.CreatedAt = time.Now()
	project.ID = primitive.NewObjectID()
	bProject, err := bson.Marshal(project)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return project.GetCollection().InsertOne(ctx, bProject)
}
