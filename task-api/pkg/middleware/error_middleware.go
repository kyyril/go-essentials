package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/models"
)

// ErrorHandlingMiddleware handles panics and logs errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error:   "internal_error",
					Message: "An internal server error occurred",
					Code:    http.StatusInternalServerError,
				})
			}
		}()
		c.Next()
	}
}

// RequestLoggingMiddleware logs incoming requests
func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
		log.Printf("%s %s - Status: %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

// CORSMiddleware configures CORS
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
