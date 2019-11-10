package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAccountsParam struct {
	RoleName string `json:"roleName"`
}

func getRoles(c *gin.Context) {
	sp := getAccountsParam{}
	if err := c.ShouldBindQuery(&sp); err == nil {
		// // revenueList := getSalesHistory(sp.Category, sp.Year)
		// c.JSON(http.StatusOK, revenueList)
	} else {
		c.Status(http.StatusBadRequest)
	}
}
