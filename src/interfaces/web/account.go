package web

import (
	"net/http"
	"reflect"
	account "williamfeng323/mooncake-duty/src/domains/account"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BasicAccount contains basic user info and validated data.
type BasicAccount struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"isAdmin"`
}

// AccountID describe the uri parameter for account
type AccountID struct {
	ID string `uri:"id" binding:"required"`
}

// AccountFilter list filter
type AccountFilter struct {
	Email string `form:"email" binding:"required,email"`
}

// AccountUpdate define the parameters to update account
type AccountUpdate struct {
	Mobile         string                 `json:"mobile,omitempty"`
	Email          string                 `json:"email,omitempty" binding:"omitempty,email"`
	ContactMethods account.ContactMethods `json:"contactMethods,omitempty"`
}

//RegisterAccountRoute the account APIs to root.
func RegisterAccountRoute(router *gin.Engine) {
	accountRoutes := router.Group("/accounts")
	{
		accountRoutes.GET("", func(c *gin.Context) {
			af := &AccountFilter{}
			if err := c.ShouldBindQuery(af); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			var acct *account.Account
			acct, err := account.GetAccountService().GetAccountByEmail(af.Email)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"account": acct})
			return
		})
		accountRoutes.PUT("", func(c *gin.Context) {
			ba := &BasicAccount{}
			if err := c.ShouldBindJSON(ba); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			token, acct, err := account.GetAccountService().Register(ba.Email, ba.Password, ba.IsAdmin)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, gin.H{"token": token, "account": acct})
			return
		})
		accountRoutes.GET("/:id", func(c *gin.Context) {
			id := &AccountID{}
			if err := c.ShouldBindUri(id); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if _, err := primitive.ObjectIDFromHex(id.ID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			acct, err := account.GetAccountService().GetAccountByID(id.ID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"account": acct})
			return
		})
		accountRoutes.POST("/:id", func(c *gin.Context) {
			id := &AccountID{}
			if err := c.ShouldBindUri(id); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			updateData := AccountUpdate{}
			zeroValue := reflect.Zero(reflect.TypeOf(updateData)).Interface()
			if err := c.ShouldBindJSON(&updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else if reflect.DeepEqual(updateData, zeroValue) {
				c.JSON(http.StatusNotModified, gin.H{})
				return
			}
			err := account.GetAccountService().UpdateContactMethods(id.ID, updateData.ContactMethods, updateData.Email, updateData.Mobile)
			if err != nil {
				if _, ok := err.(account.NotFoundError); ok {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
			return
		})
	}
	router.POST("/login", func(c *gin.Context) {
		ba := &BasicAccount{}
		if err := c.ShouldBindJSON(ba); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := account.GetAccountService().SignIn(ba.Email, ba.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	})
	router.POST("/refresh", func(c *gin.Context) {
		type tokenStruct struct {
			Token string `json:"token" bind:"required"`
		}
		tk := &tokenStruct{}
		if err := c.ShouldBindJSON(tk); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		refreshedTk, err := account.GetAccountService().Refresh(tk.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token validation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": refreshedTk})
		return
	})
}
