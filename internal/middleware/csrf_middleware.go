package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware basic CSRF protection for cookie-based auth
// It checks if the Origin or Referer header matches the allowed origins
// and requires a custom header for state-changing requests.
func CSRFMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only check for state-changing methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 1. Check Origin header
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = c.GetHeader("Referer")
		}

		isAllowed := false
		for _, o := range allowedOrigins {
			if strings.HasPrefix(origin, o) {
				isAllowed = true
				break
			}
		}

		if !isAllowed && origin != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF: Origin not allowed"})
			c.Abort()
			return
		}

		// 2. Custom header check (Double Submit Cookie variant or just specific header)
		// For simplicity, we require a custom header like "X-Requested-With" or "X-CSRF-Token"
		// which cannot be sent cross-origin without CORS approval and script access.
		if c.GetHeader("X-Requested-With") == "" && c.GetHeader("X-CSRF-Token") == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF: Missing required security headers"})
			c.Abort()
			return
		}

		c.Next()
	}
}
