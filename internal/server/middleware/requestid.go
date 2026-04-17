package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := make([]byte, 12)
		_, _ = rand.Read(buf)
		requestID := hex.EncodeToString(buf)
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}
