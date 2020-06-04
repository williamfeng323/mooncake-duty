package project

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	db "williamfeng323/mooncake-duty/src/infrastructure/db"
)

// Permission the predefined permissions strings
type Permission string

// Tier is the support level of a team member.
type Tier string

// Severity describe the level of the alarm.
type Severity int

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

const (
	// High represents the most serious issue, must be handle asap
	High Severity = iota
	// Medium represents the medium level issue, you might not need to escalate
	Medium
	// Low could be follow up later
	Low
)

// Alarm defines the struct of an alarm
type Alarm struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Severity    `json:"severity" bson:"severity"`
}

// Project is the struct to contain project information.
type Project struct {
	db.BaseModel  `json:",inline" bson:",inline"`
	Name          string               `json:"name" bson:"name" required:"true"`
	Description   string               `json:"description" bson:"description" required:"true"`
	Members       []primitive.ObjectID `json:"members" bson:"members"`
	ProjectAdmins []primitive.ObjectID `json:"projectAdmins" bson:"projectAdmins"`
	Alarms        []Alarm              `json:"alarms" bson:"alarms"` // AlarmLog would be needed
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
