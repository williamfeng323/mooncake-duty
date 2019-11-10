package account

import (
	"context"
	"time"
	"williamfeng323/mooncake-duty/src/dao"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Account struct of the user account
type Account struct {
	dao.BaseModel `json:",inline" bson:",inline"`
	Projects      []string `json:"projects" bson:"projects"`
	Email         string   `json:"email" bson:"email" required:"true"`
	Password      string   `json:"password" bson:"password" required:"true"`
}

// InsertAccount create a new document in mongoDB with
// initial createdAt and _id
func (acct *Account) InsertAccount() (*mongo.InsertOneResult, error) {
	validationErrors := acct.DefaultValidator()
	if validationErrors != nil {
		return nil, validationErrors[0]
	}
	acct.CreatedAt = time.Now()
	acct.ID = primitive.NewObjectID()
	bAcct, err := bson.Marshal(acct)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return acct.GetCollection().InsertOne(ctx, bAcct)
}
