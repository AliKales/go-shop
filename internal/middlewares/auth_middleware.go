package middlewares

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTokenFromHeader(c *gin.Context) *string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil
	}

	token := authHeader[len(bearerPrefix):]

	return &token
}

func UserAuthMiddleware(c *gin.Context) {
	if c.IsAborted() {
		return
	}

	// Get the Bearer token from the "Authorization" header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
		c.Abort()
		return
	}

	// Check if the Authorization header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	// Extract the token
	token := authHeader[len(bearerPrefix):]

	user := database.GetUserBy("token", &token)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		c.Abort()
		return
	}

	if utils.IsExpired(user.TokenExpireAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token expired!"})
		c.Abort()
		return
	}

	c.Set("user", user)

	c.Next()
}

func ConvertToUser(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}

	authUser, ok := user.(*models.User)
	if !ok {
		return nil
	}

	return authUser
}
