package models

import (
	tokengenerator "example/web-service-gin/internal/token_generator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/pgtype"
)

type User struct {
	Id                   uint      `gorm:"primaryKey"`
	Email                string    `gorm:"type:varchar(255);not null;unique"`
	Token                string    `gorm:"type:varchar(255);not null"`
	TokenExpireAt        time.Time `gorm:"not null"`
	RefreshToken         string    `gorm:"type:varchar(255);not null"`
	RefreshTokenExpireAt time.Time `gorm:"not null"`
	Username             string    `gorm:"type:varchar(40);unique"`
	Password             string    `gorm:"type:varchar(255);not null"`
	CreatedAt            time.Time `gorm:"not null"`
	IsEmailVerified      bool      `gorm:"default:false"`
	UserAction           string    `gorm:"type:text"`
	UserActionExpireAt   time.Time `gorm:"not null"`
	StoreId              *int
	Cart                 *pgtype.JSONB
}

func (u *User) ResetAllTokens() {
	u.ResetToken()
	u.RefreshToken = tokengenerator.GenerateUserRefreshToken()
	u.RefreshTokenExpireAt = time.Now().Add(24 * time.Hour).UTC()
}

func (u *User) ResetToken() {
	u.Token = tokengenerator.GenerateUserToken()
	u.TokenExpireAt = time.Now().Add(30 * time.Minute).UTC()
}

func (u *User) TokenData() gin.H {
	return gin.H{
		"token":                u.Token,
		"tokenExpireAt":        u.TokenExpireAt,
		"refreshToken":         u.RefreshToken,
		"refreshTokenExpireAt": u.RefreshTokenExpireAt,
	}
}

func (u *User) PublicData() gin.H {
	return gin.H{
		"username":  u.Username,
		"createdAt": u.CreatedAt,
	}
}

func (u *User) PrivateData() gin.H {
	return gin.H{
		"username":        u.Username,
		"createdAt":       u.CreatedAt,
		"id":              u.Id,
		"email":           u.Email,
		"isEmailVerified": u.IsEmailVerified,
		"storeId":         u.StoreId,
		"cartItemCount": u.CartItemCount(),
	}
}

func (u *User) CartItemCount() int {
	var data map[string]int
	u.Cart.AssignTo(&data)

	return len(data)
}
