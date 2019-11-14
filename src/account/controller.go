package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, createResult)
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
