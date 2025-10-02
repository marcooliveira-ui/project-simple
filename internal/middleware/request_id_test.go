package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should generate new request ID if not provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		handler := RequestID()
		handler(c)

		requestID := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, requestID)

		// Verify it's a valid UUID
		_, err := uuid.Parse(requestID)
		assert.NoError(t, err)

		// Verify it's stored in context
		ctxRequestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.Equal(t, requestID, ctxRequestID)
	})

	t.Run("Should use existing request ID from header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		existingID := uuid.New().String()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", existingID)
		c.Request = req

		handler := RequestID()
		handler(c)

		requestID := w.Header().Get("X-Request-ID")
		assert.Equal(t, existingID, requestID)

		ctxRequestID, exists := c.Get("request_id")
		assert.True(t, exists)
		assert.Equal(t, existingID, ctxRequestID)
	})

	t.Run("Should set response header with request ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		handler := RequestID()
		handler(c)

		assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	})
}
