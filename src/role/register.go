package role

import (
	"github.com/gin-gonic/gin"
)

//Register the sales APIs to root.
func Register(router *gin.RouterGroup) {
	saleRoutes := router.Group("/roles")
	{
		saleRoutes.GET("", getRoles)
	}
}
