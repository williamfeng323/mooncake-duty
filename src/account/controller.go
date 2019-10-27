package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type saleParam struct {
	Category string `form:"category"`
	Year     int16  `form:"year"`
}

func getSale(c *gin.Context) {
	sp := saleParam{}
	if err := c.ShouldBindQuery(&sp); err == nil {
		// // revenueList := getSalesHistory(sp.Category, sp.Year)
		// c.JSON(http.StatusOK, revenueList)
	} else {
		c.Status(http.StatusBadRequest)
	}
}
