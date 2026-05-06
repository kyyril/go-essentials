package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"task-management-api/pkg/config"
	"task-management-api/pkg/models"
)

// GenerateTokens generates access and refresh tokens
func GenerateTokens(user *models.User, cfg *config.Config) (accessToken, refreshToken string, err error) {
	accessToken, err = generateToken(user, "access", cfg.JWT.ExpirationSeconds, cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = generateToken(user, "refresh", cfg.JWT.RefreshExpSeconds, cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken generates a JWT token
func generateToken(user *models.User, tokenType string, expirationSeconds int, secret string) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"type":  tokenType,
		"exp":   time.Now().Add(time.Duration(expirationSeconds) * time.Second).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ExtractUserIDFromToken extracts user ID from JWT claims
func ExtractUserIDFromToken(claims jwt.MapClaims) (string, error) {
	userID, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID in token")
	}
	return userID, nil
}

// ExtractRoleFromToken extracts role from JWT claims
func ExtractRoleFromToken(claims jwt.MapClaims) (string, error) {
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("invalid role in token")
	}
	return role, nil
}

// IsTokenExpired checks if token is expired
func IsTokenExpired(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return true
	}
	return time.Now().Unix() > int64(exp)
}
