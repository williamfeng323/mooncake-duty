package account

import (
	"context"
	"time"
	"williamfeng323/mooncake-duty/src/dao"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SentHook is the structure to describe the alternative way to send the alarms.
type SentHook struct {
	URL string `json:"url" bson:"url"`
}

// ContactMethods is describing the way to contact user
type ContactMethods struct {
	SentHook  `json:"sentHook,inline" bson:"sentHook,inline"`
	SentSMS   bool `json:"sentSMS,omitempty" bson:"sentSMS,omitempty"`
	SentEmail bool `json:"sentEmail,omitempty" bson:"sentEmail,omitempty"`
	IsEnabled bool `json:"isEnabled" bson:"isEnabled"`
}

// Account struct of the user account
type Account struct {
	dao.BaseModel  `json:",inline" bson:",inline"`
	Email          string               `json:"email" bson:"email" required:"true"`
	Password       string               `json:"password" bson:"password" required:"true"`
	Mobile         string               `json:"mobile,omitempty" bson:"mobile,omitempty"`
	IsAdmin        bool                 `json:"isAdmin" bson:"isAdmin"`
	Avatar         string               `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Teams          []primitive.ObjectID `json:"teams,omitempty" bson:"teams, omitempty"`
	ContactMethods `json:"contactMethods,omitempty" bson:"contactMethods,omitempty"`
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
