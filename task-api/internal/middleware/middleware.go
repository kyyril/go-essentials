package middleware

import (
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("thisIsSecretKeyENV")

func AuthMiddleware(next http.Handler, requiredRole string) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer"){
			http.Error(w, "Unauthorizad", http.StatusUnauthorized)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer")
		claims := jwt.MapClaims{}

		//validasi token
		_, err := jwt.ParseWithClaims(tokenStr, claims, func (token *jwt.Token) (interface{}, error)  {
			return []byte(jwtKey), nil
		})
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		//check role
		if requiredRole != "" && claims["role"] != requiredRole {
			http.Error(w, "Forbidden: You don't have access", http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}