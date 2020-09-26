package middlewares

import (
	"net/http"
	"williamfeng323/mooncake-duty/src/domains/account"
	"williamfeng323/mooncake-duty/src/utils"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// AuthHeader defines the header of the authentication
const AuthHeader string = "Authorization"

// Authenticate middleware to check the request is authorized.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(AuthHeader)
		if len(token) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": jwt.ValidationError{}})
			return
		}
		claims, err := utils.VerifyToken(token, false)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err})
			return
		}
		acct, err := account.GetAccountService().GetAccountByEmail(claims.Audience)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err})
			return
		}
		c.Set("user", acct)
		c.Next()
	}
}
