package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request for tracing
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get request ID from header
		requestID := c.GetHeader("X-Request-ID")

		// Generate new ID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for handlers to access
		c.Set("request_id", requestID)

		// Add to response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
