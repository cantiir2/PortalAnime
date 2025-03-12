package middleware

import (
	"net/http"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/services"
)

// AuthMiddleware creates a middleware for authentication
func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		userID, role, err := userService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Next()
	}
}

// AdminMiddleware creates a middleware for admin-only routes
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// Tambahkan debug logging
		log.Printf("Role from context: %v (type: %T)", roleInterface, roleInterface)

		role, ok := roleInterface.(models.Role)
		if !ok {
			// Coba konversi dari string jika bukan Role type
			roleStr, ok := roleInterface.(string)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
				c.Abort()
				return
			}
			role = models.Role(roleStr)
		}

		// Tambahkan debug logging
		log.Printf("Checking admin access - Role: %v, Expected: %v, Equal: %v",
			role, models.RoleAdmin, role == models.RoleAdmin)

		if role != models.RoleAdmin {
			// Ubah pesan error menjadi lebih deskriptif
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required ASDASDA"}) // Hapus "ASDASDAS"
			c.Abort()
			return
		}

		c.Next()
	}
}
