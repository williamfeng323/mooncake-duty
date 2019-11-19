package account

import (
	"github.com/gin-gonic/gin"
)

//Register the sales APIs to root.
func Register(router *gin.RouterGroup) {
	saleRoutes := router.Group("/accounts")
	{
		saleRoutes.PUT("", createAccountController)
		saleRoutes.POST("/:id", getAccountByIDController)
		saleRoutes.POST("/login", loginController)
		saleRoutes.POST("/refresh", refreshController)
	}
}
