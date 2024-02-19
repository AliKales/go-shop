package models

import "time"

type StoreItem struct {
	Id            uint      `json:"id"`
	StoreId       int       `json:"storeId"`
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	ItemCount     int       `json:"itemCount"`
	ItemSellCount int       `json:"itemSellCount"`
	Price         float64   `json:"price"`
}
