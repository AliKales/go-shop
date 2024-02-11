package handlers

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	tokengenerator "example/web-service-gin/internal/token_generator"
	"example/web-service-gin/internal/utils"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SignupReq struct {
	Email    string
	Password string
	Username string
}

type LoginReq struct {
	Password string
	Username string
}

type ChangePassReq struct{
	Password string
	Token string
}

func SignupHandler(c *gin.Context) {
	var req SignupReq
	c.BindJSON(&req)

	hash_pass := utils.HashPassword(req.Password)

	action := tokengenerator.GenerateSecureToken(20) + "|verify-email"

	user := models.User{Username: req.Username, Email: req.Email, Password: hash_pass, Token: tokengenerator.GenerateUserToken(), RefreshToken: tokengenerator.GenerateUserRefreshToken(), TokenExpireAt: time.Now().Add(30 * time.Minute).UTC(), RefreshTokenExpireAt: time.Now().Add(24 * time.Hour).UTC(), CreatedAt: time.Now(), IsEmailVerified: false, UserAction: action, UserActionExpireAt: time.Now().Add(100 * time.Hour).UTC()}

	if err := database.DB.Create(&user).Error; err != nil {
		if utils.IsNotUniqueColumn(err, "users_email_key") {
			c.JSON(http.StatusConflict, gin.H{"message": "Email is already taken", "val": "email"})
			return
		} else if utils.IsNotUniqueColumn(err, "users_username_key") {
			c.JSON(http.StatusConflict, gin.H{"message": "Username is already taken", "val": "username"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	utils.SendEmail(user.Email, "Verify email", "To verify your email please click this; https://storei-rp46.onrender.com/auth/verify-email?token="+action)

	c.JSON(http.StatusForbidden, gin.H{"message": "You are signed up! Now verify your email"})
}

func LoginHandler(c *gin.Context) {
	var req LoginReq
	c.BindJSON(&req)

	user := database.GetUserBy("username", &req.Username)

	if user == nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Username and password is wrong"})
		return
	}

	if !user.IsEmailVerified {
		if utils.IsExpired(user.UserActionExpireAt) {
			newAction := tokengenerator.GenerateSecureToken(20) + "|verify-email"
			user.UserAction = newAction
			user.UserActionExpireAt = time.Now().Add(100 * time.Minute).UTC()
			database.DB.Save(&user)

			utils.SendEmail(user.Email, "Verify email", "To verify your email please click this; https://storei-rp46.onrender.com/auth/verify-email?token="+newAction)
		}
		c.JSON(http.StatusForbidden, gin.H{"message": "Verify your email"})
		return
	}

	if utils.IsExpired(user.RefreshTokenExpireAt) {
		user.ResetAllTokens()
		database.DB.Save(&user)
	} else if utils.IsExpired(user.TokenExpireAt) {
		user.ResetToken()
		database.DB.Save(&user)
	}

	c.JSON(http.StatusOK, user.TokenData())
}

func RefreshTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
		return
	}

	// Check if the Authorization header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
		return
	}

	// Extract the token
	refreshToken := authHeader[len(bearerPrefix):]

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token is not present!"})
		return
	}

	user := database.GetUserBy("refresh_token", &refreshToken)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found!"})
		return
	}

	if utils.IsExpired(user.RefreshTokenExpireAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token expired. Please re-login!"})
		return
	}

	if !utils.IsExpired(user.TokenExpireAt) {
		c.JSON(http.StatusOK, user.TokenData())
		return
	}

	user.ResetToken()

	database.DB.Save(&user)

	c.JSON(http.StatusOK, user.TokenData())
}

func VerifyEmailHandler(c *gin.Context, user models.User) {
	user.IsEmailVerified = true
	user.UserAction = ""
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Your email is verified!"})
}

func ForgetPasswordHandler(c *gin.Context) {
	email := c.Query("email")

	user := database.GetUserBy("email", &email)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found!"})
		return
	}

	if user.UserAction != "" && !utils.IsExpired(user.UserActionExpireAt) && strings.Contains(user.UserAction, "change-password") {
		diff := user.UserActionExpireAt.Sub(time.Now().UTC())

		minute := fmt.Sprintf("%.0f", math.Floor(diff.Minutes()))

		c.JSON(http.StatusTooManyRequests, gin.H{"message": "We already sent you email. You can request again in " + minute + " minutes!", "val": minute})
		return
	}

	newAction := tokengenerator.GenerateSecureToken(20) + "|change-password"
	user.UserAction = newAction
	user.UserActionExpireAt = time.Now().Add(20 * time.Minute).UTC()
	database.DB.Save(&user)

	utils.SendEmail(user.Email, "Reset password", "To reset your password please click this; https://storei-rp46.onrender.com/auth/reset-password?token="+newAction)

	c.JSON(http.StatusOK, gin.H{"message": "We sent you an email to reset your password!"})
}

func ChangePasswordHandler(c *gin.Context) {
	var req ChangePassReq
	c.BindJSON(&req)

	user := database.GetUserBy("user_action", &req.Token)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Token does not exist!"})
		return
	}

	if utils.IsExpired(user.UserActionExpireAt) {
		c.JSON(http.StatusGone, gin.H{"message": "Token expired!"})
		return
	}

	action := strings.Split(user.UserAction, "|")[1]
	if action != "change-password" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error!"})
		return
	}

	user.Password = utils.HashPassword(req.Token)
	user.UserAction = ""
	user.ResetAllTokens()

	database.DB.Save(user)

	c.JSON(http.StatusOK, gin.H{"message": "Your password has changed!"})
}

func RequestDeleteAccountHandler(c *gin.Context) {
	var req LoginReq
	c.BindJSON(&req)

	user := database.GetUserBy("username", &req.Username)

	if user == nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Username and password is wrong"})
		return
	}

	if user.UserAction != "" && !utils.IsExpired(user.UserActionExpireAt) && strings.Contains(user.UserAction, "delete-account") {
		diff := user.UserActionExpireAt.Sub(time.Now().UTC())

		minute := fmt.Sprintf("%.0f", math.Floor(diff.Minutes()))

		c.JSON(http.StatusTooManyRequests, gin.H{"message": "We already sent you email. You can request again in " + minute + " minutes!", "val": minute})
		return
	}

	newAction := tokengenerator.GenerateSecureToken(20) + "|delete-account"
	user.UserAction = newAction
	user.UserActionExpireAt = time.Now().Add(10 * time.Minute).UTC()
	database.DB.Save(&user)

	utils.SendEmail(user.Email, "Delete Account", "To delete your account please click this; https://storei-rp46.onrender.com/auth/delete-account?token="+newAction)

	c.JSON(http.StatusOK, gin.H{"message": "We sent you an email to delete your account!"})
}

func DeleteAccountHandler(c *gin.Context, user models.User) {
	database.DB.Delete(user)

	store := database.GetStoreBy("user_id", user.Id)

	if store == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User deleted!"})
		return
	}

	database.DB.Delete(store)

	database.DB.Where("store_id = ?", store.ID).Delete(&models.Store{})

	c.JSON(http.StatusOK, gin.H{"message": "User deleted!"})
}
