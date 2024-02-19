package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Store struct {
	Id            uint
	UserId        int
	CreatedAt     time.Time
	Name          string
	LinkName      string
	ItemCount     int
	ItemSellCount int
}

func (s *Store) PublicData() gin.H {
	return gin.H{
		"name":      s.Name,
		"linkName":  s.LinkName,
		"createdAt": s.CreatedAt,
		"itemCount": s.ItemCount,
		"itemSellCount": s.ItemSellCount,
		"userId":    s.UserId,
	}
}