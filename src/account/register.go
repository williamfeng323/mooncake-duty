package account

import (
	"github.com/gin-gonic/gin"
)

//Register the account APIs to root.
func Register(router *gin.Engine) {
	accountRoutes := router.Group("/accounts")
	{
		accountRoutes.PUT("", createAccountController)
		accountRoutes.GET("/:id", getAccountByIDController)
		accountRoutes.POST("/login", loginController)
		accountRoutes.POST("/refresh", refreshController)
	}
}
