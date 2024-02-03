package models

import (
	"time"

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
	Cart                 pgtype.JSONB
}
