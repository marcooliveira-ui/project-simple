package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func BadRequest(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "Bad Request",
		Message: message,
		Details: details,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "Not Found",
		Message: message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "Internal Server Error",
		Message: message,
	})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Error:   "Conflict",
		Message: message,
	})
}

func UnprocessableEntity(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Error:   "Unprocessable Entity",
		Message: message,
		Details: details,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Message: message,
		Data:    data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
