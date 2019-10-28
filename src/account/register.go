package account

import (
	"github.com/gin-gonic/gin"
)

//Register the sales APIs to root.
func Register(router *gin.RouterGroup) {
	saleRoutes := router.Group("/accounts")
	{
		saleRoutes.GET("", getAccounts)
	}
}
