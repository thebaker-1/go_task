package Infrastructure

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key should be stored securely, e.g., environment variable
// Use package-level jwtSecret from jwt_services.go


func AuthenticateJWT() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Log the unexpected method for debugging, but don't expose too much info to client
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		// CRITICAL: Handle the error from jwt.Parse
		if err != nil {
			log.Printf("Error parsing token: %v", err) // Log the actual error for server-side debugging
			// Provide a generic error message to the client for security
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user information in context for subsequent handlers
			// Type assertion for safety
			// if userID, ok := claims["user_id"].(float64); ok { // JWT numbers are float64
			// 	c.Set("user_id", int(userID)) // Store as int if appropriate
			// } else {
			// 	log.Printf("user_id claim is not a number: %v", claims["user_id"])
			// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			// 	c.Abort()
			// 	return
			// }

			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			} else {
				log.Printf("username claim is not a string: %v", claims["username"])
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username in token"})
				c.Abort()
				return
			}

			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			} else {
				log.Printf("role claim is not a string: %v", claims["role"])
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
				c.Abort()
				return
			}

			c.Next()
		} else {
			// This block should ideally be rarely hit if jwt.Parse error handling is robust,
			// but it catches cases where token is not valid *after* parsing (e.g., claims issues).
			log.Printf("Token is not valid or claims are invalid: %v", token.Valid)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Access denied. Requires '%s' role.", requiredRole)})
			c.Abort()
			return
		}
		c.Next()
	}
}
