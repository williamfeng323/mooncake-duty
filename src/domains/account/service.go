package account

import (
	"context"
	"fmt"
	"time"

	db "williamfeng323/mooncake-duty/src/infrastructure/db"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type basicAccountsParam struct {
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// 	IsAdmin  bool   `json:"isAdmin" binding:"required"`
// }

// func CreateAccountController(c *gin.Context) {
// 	sp := basicAccountsParam{}
// 	if err := c.ShouldBind(&sp); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	createResult, err := createAccount(sp.Email, sp.Password, sp.IsAdmin)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	id := createResult.InsertedID.(primitive.ObjectID)
// 	c.JSON(http.StatusOK, gin.H{"id": id.String()})
// }

// type UpdateAccountParam struct {
// 	Avatar string `json:"email"`
// 	Mobile string `json:"mobile"`
// }

// func UpdateAccountController(c *gin.Context) {
// 	var id string
// 	sp := UpdateAccountParam{}

// 	if err := c.ShouldBindUri(&id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}
// 	if err := c.ShouldBind(&sp); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	convertedID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}
// 	rst, err := updateAccount(convertedID, sp.Avatar, sp.Mobile)
// 	if err != nil {
// 		c.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"modifiedCount": rst.ModifiedCount})
// }

// func GetAccountByIDController(c *gin.Context) {
// 	var id string
// 	if err := c.ShouldBindUri(&id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}
// 	user, err := getAccountByID(id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// func LoginController(c *gin.Context) {
// 	sp := basicAccountsParam{}
// 	if err := c.ShouldBind(&sp); err != nil {
// 		c.Status(http.StatusBadRequest)
// 	}
// 	tokenString, err := signIn(sp.Email, sp.Password)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// 	return
// }

// type RefreshParam struct {
// 	Token string `json:"token"`
// }

// func RefreshController(c *gin.Context) {
// 	params := RefreshParam{}
// 	if err := c.ShouldBind(&params); err != nil {
// 		c.Status(http.StatusBadRequest)
// 		return
// 	}
// 	refreshedToken, err := refresh(params.Token)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"token": refreshedToken})
// }

// Service exposes the account service
type Service struct {
	repo db.Repository
}

// SetRepo set the account repository to the service
func (as *Service) SetRepo(repo *repoimpl.AccountRepo) {
	as.repo = repo
}

// SignIn validate the email password pair and return token
func (as *Service) SignIn(email string, password string) (string, error) {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", utils.GetConf().Mongo.DefaultTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	result := as.repo.FindOne(ctx, bson.M{"email": email})
	acct := &Account{}
	result.Decode(acct)
	if acct.ID.IsZero() {
		return "", fmt.Errorf("Account does not exist")
	}
	decryptedPassword, err := utils.Decrypt(acct.Password)
	if err != nil {
		return "", err
	}
	if string(decryptedPassword) != password {
		return "", fmt.Errorf("Password and account does not match")
	}
	token, err := utils.SignToken(acct.Email)
	return token, err
}
