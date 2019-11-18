package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type basicAccountsParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createAccountController(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.Status(http.StatusInternalServerError)
		}
	}()
	sp := basicAccountsParam{}
	if err := c.ShouldBind(&sp); err != nil {
		c.Status(http.StatusBadRequest)
	}
	createResult, err := createAccount(sp.Email, sp.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	id := createResult.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func loginController(c *gin.Context) {
	sp := basicAccountsParam{}
	if err := c.ShouldBind(&sp); err != nil {
		c.Status(http.StatusBadRequest)
	}
	tokenString, err := signIn(sp.Email, sp.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
	return
}

type refreshParam struct {
	Token string `json:"token"`
}

func refreshController(c *gin.Context) {
	params := refreshParam{}
	if err := c.ShouldBind(&params); err != nil {
		c.Status(http.StatusBadRequest)
	}
	refreshedToken, err := refresh(params.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"token": refreshedToken})
	return
}
