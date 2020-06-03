package web

import (
	"github.com/gin-gonic/gin"

	account "williamfeng323/mooncake-duty/src/domains/account"
)

//RegisterAccountRoute the account APIs to root.
func RegisterAccountRoute(router *gin.Engine) {
	accountRoutes := router.Group("/accounts")
	{
		accountRoutes.PUT("", account.CreateAccountController)
		accountRoutes.GET("/:id", account.GetAccountByIDController)
		accountRoutes.POST("/login", account.LoginController)
		accountRoutes.POST("/refresh", account.RefreshController)
	}
}
