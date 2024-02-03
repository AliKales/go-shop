package controllers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func ResetAllTokens(user models.User) *models.User {
	user.Token = utils.GenerateSecureToken(16)
	user.TokenExpireAt = time.Now().Add(10 * time.Minute).UTC()
	user.RefreshToken = utils.GenerateSecureToken(32)
	user.RefreshTokenExpireAt = time.Now().Add(30 * time.Minute).UTC()
	return &user
}

func ResetToken(user models.User) *models.User {
	user.Token = utils.GenerateSecureToken(16)
	user.TokenExpireAt = time.Now().Add(10 * time.Minute).UTC()
	return &user
}

func UserTokenData(user models.User) *gin.H {
	return &gin.H{
		"token":                user.Token,
		"tokenExpireAt":        user.TokenExpireAt,
		"refreshToken":         user.RefreshToken,
		"refreshTokenExpireAt": user.RefreshTokenExpireAt,
	}
}

func UserPublicData(user *models.User) *gin.H {
	return &gin.H{
		"username":  user.Username,
		"createdAt": user.CreatedAt,
	}
}

func UserCartItemCount(user models.User) int {
	var data map[string]int
	user.Cart.AssignTo(&data)

	return len(data)
}
