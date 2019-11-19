package account

import (
	"context"
	"fmt"
	"time"
	"williamfeng323/mooncake-duty/src/dao"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createAccount(email string, password string) (*mongo.InsertOneResult, error) {
	acct := &Account{}
	conn := dao.GetConnection()
	acctModel := conn.GetCollection("Account")
	acctModel.New(acct)
	acct.Email = email
	acct.Password = password
	isExist := acctModel.FindOne(nil, bson.M{"email": email})
	if isExist.Err() == nil {
		return nil, fmt.Errorf("Account already exists")
	}
	return acct.InsertAccount()
}

func getAccountByID(id string) (Account, error) {
	var acct Account
	conn := dao.GetConnection()
	acctModel := conn.GetCollection("Account")
	acctModel.New(&acct)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Account{}, err
	}
	result := acct.FindByID(objID)
	err = result.Decode(&acct)
	if err != nil {
		return Account{}, err
	}
	return acct, nil
}

func signIn(email string, password string) (string, error) {
	acct := &Account{}
	conn := dao.GetConnection()
	acctModel := conn.GetCollection("Account")
	acctModel.New(acct)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	cursor := acctModel.FindOne(ctx, bson.M{
		"email":    email,
		"password": password,
	})
	err := cursor.Decode(acct)
	if err != nil {
		return "", err
	}
	jwt, err := utils.SignToken(acct.Email)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func refresh(tokenString string) (string, error) {
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	refreshedToken, err := utils.SignToken(claims.Audience)
	if err != nil {
		return "", err
	}
	return refreshedToken, nil
}
