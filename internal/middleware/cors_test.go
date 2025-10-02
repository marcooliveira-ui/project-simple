package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Request without Origin header should be allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		handler := CORS([]string{"http://localhost:3000"})
		handler(c)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.False(t, c.IsAborted())
	})

	t.Run("Request with allowed Origin should be accepted", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		c.Request = req

		handler := CORS([]string{"http://localhost:3000", "http://localhost:8080"})
		handler(c)

		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.False(t, c.IsAborted())
	})

	t.Run("Request with disallowed Origin should be rejected", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://malicious-site.com")
		c.Request = req

		handler := CORS([]string{"http://localhost:3000"})
		handler(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Wildcard origin should allow all", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://any-site.com")
		c.Request = req

		handler := CORS([]string{"*"})
		handler(c)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.False(t, c.IsAborted())
	})

	t.Run("OPTIONS request should return 204", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		c.Request = req

		handler := CORS([]string{"http://localhost:3000"})
		handler(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Should set correct CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		handler := CORS([]string{"http://localhost:3000"})
		handler(c)

		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
	})
}
