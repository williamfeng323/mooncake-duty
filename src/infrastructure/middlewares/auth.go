package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Auth middleware to check the request is authorized.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// before request
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		rawRequestParam, _ := json.Marshal(c.Request.URL.Query())
		log.Printf("\"Request Params\": %s", rawRequestParam)

		c.Next()

		// after request
		latency := time.Since(t)
		statusCode := c.Writer.Status()
		fmt.Printf("{\"responseBody\": %s, \"status\": %d, \"responseTime\": %f}", blw.body.String(), statusCode, latency.Seconds())
	}
}
