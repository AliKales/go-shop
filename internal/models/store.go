package models

import (
	"time"
)

type Store struct {
	ID            uint
	UserId        int
	CreatedAt     time.Time
	Name          string
	LinkName      string
	ItemCount     int
	ItemSellCount int
}
