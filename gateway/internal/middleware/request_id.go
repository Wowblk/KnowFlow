package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestID ensures every request has a stable correlation id across gateway,
// backend, and agent services.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
			c.Request.Header.Set(requestIDHeader, requestID)
		}
		c.Header(requestIDHeader, requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}
