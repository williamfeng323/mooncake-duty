package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger middleware to log request before and after.
func Logger() gin.HandlerFunc {
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
