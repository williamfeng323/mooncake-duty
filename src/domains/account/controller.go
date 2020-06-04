package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type basicAccountsParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"isAdmin" binding:"required"`
}

func CreateAccountController(c *gin.Context) {
	sp := basicAccountsParam{}
	if err := c.ShouldBind(&sp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createResult, err := createAccount(sp.Email, sp.Password, sp.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := createResult.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

type UpdateAccountParam struct {
	Avatar string `json:"email"`
	Mobile string `json:"mobile"`
}

func UpdateAccountController(c *gin.Context) {
	var id string
	sp := UpdateAccountParam{}

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := c.ShouldBind(&sp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	convertedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	rst, err := updateAccount(convertedID, sp.Avatar, sp.Mobile)
	if err != nil {
		c.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"modifiedCount": rst.ModifiedCount})
}

func GetAccountByIDController(c *gin.Context) {
	var id string
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user, err := getAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}

func LoginController(c *gin.Context) {
	sp := basicAccountsParam{}
	if err := c.ShouldBind(&sp); err != nil {
		c.Status(http.StatusBadRequest)
	}
	tokenString, err := signIn(sp.Email, sp.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
	return
}

type RefreshParam struct {
	Token string `json:"token"`
}

func RefreshController(c *gin.Context) {
	params := RefreshParam{}
	if err := c.ShouldBind(&params); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	refreshedToken, err := refresh(params.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": refreshedToken})
}