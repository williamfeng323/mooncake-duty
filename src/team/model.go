package team

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

// Team is the struct to contain team information.
type Team struct {
	dao.BaseModel `json:",inline" bson:",inline"`
	Name          string   `json:"name" bson:"name" required:"true"`
	Description   string   `json:"description" bson:"description" required:"true"`
	Members       []Member `json:"members" bson:"members"`
}

// InsertTeam create a new team document in mongoDB with
// initial createdAt and _id
func (team *Team) InsertTeam() (*mongo.InsertOneResult, error) {
	validationErrors := team.DefaultValidator()
	if validationErrors != nil {
		return nil, validationErrors[0]
	}
	team.CreatedAt = time.Now()
	team.ID = primitive.NewObjectID()
	bTeam, err := bson.Marshal(team)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return team.GetCollection().InsertOne(ctx, bTeam)
}
