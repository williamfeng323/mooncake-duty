package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAccountsParam struct {
	RoleName string `json:"roleName"`
}

func GetTeams(c *gin.Context) {
	sp := GetAccountsParam{}
	if err := c.ShouldBindQuery(&sp); err == nil {
		// // revenueList := getSalesHistory(sp.Category, sp.Year)
		// c.JSON(http.StatusOK, revenueList)
	} else {
		c.Status(http.StatusBadRequest)
	}
}
