package account

import (
	"context"
	"time"
	"williamfeng323/mooncake-duty/src/dao"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Role the role struct type for account.
type Role struct {
	dao.BaseModel
	Name string `json:"name" bson:"name"`
}

// Account struct of the user account
type Account struct {
	dao.BaseModel
	Projects []string             `json:"projects" bson:"projects"`
	Email    string               `json:"email" bson:"email"`
	Password string               `json:"password" bson:"password"`
	RolesID  []primitive.ObjectID `json:"rolesId" bson:"rolesId"`
}

//FindByRoleID finds the account by specific role's id
func (acct *Account) FindByRoleID(roleID primitive.ObjectID) ([]Account, error) {
	filter := bson.M{"RolesID": bson.A{roleID}}
	rst, err := acct.Find(filter)
	accts := []Account{}
	rst.All(nil, accts)
	return accts, err
}

// InsertAccount create a new document in mongoDB with
// initial createdAt and _id
func (acct *Account) InsertAccount() (*mongo.InsertOneResult, error) {
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
