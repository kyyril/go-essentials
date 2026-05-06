package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/config"
	"task-management-api/pkg/models"
	"task-management-api/pkg/utils"
)

// AuthMiddleware validates JWT tokens and injects user context
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header required",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract token from Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid authorization header format",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired token",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract user info from claims
		userID, err := utils.ExtractUserIDFromToken(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid token claims",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		role, _ := utils.ExtractRoleFromToken(claims)

		// Inject user context into request
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	return userID.(string)
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) string {
	role, exists := c.Get("user_role")
	if !exists {
		return ""
	}
	return role.(string)
}

// RequireAdmin middleware checks if user is admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRole(c)
		if role != "admin" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "forbidden",
				Message: "Admin role required",
				Code:    http.StatusForbidden,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
