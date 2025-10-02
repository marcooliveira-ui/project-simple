package middleware

import (
	"log"
	"net/http"
	"project-simple/pkg/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			log.Printf("Error: %v", err.Err)

			// If response was already written, don't write again
			if c.Writer.Written() {
				return
			}

			response.InternalServerError(c, "An unexpected error occurred")
			c.Abort()
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				c.JSON(http.StatusInternalServerError, response.ErrorResponse{
					Error:   "Internal Server Error",
					Message: "An unexpected error occurred",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
