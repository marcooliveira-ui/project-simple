package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestSizeLimit limits the maximum size of request bodies
func RequestSizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		c.Next()

		// Check if size limit was exceeded
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "Request Too Large",
				"message": "Request body exceeds maximum allowed size",
			})
			c.Abort()
		}
	}
}
