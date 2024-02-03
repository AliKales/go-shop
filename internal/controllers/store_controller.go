package controllers

import (
	"example/web-service-gin/internal/models"
	"regexp"

	"github.com/gin-gonic/gin"
)

func IsLinkNameOK(val string) bool {
	pattern := "^[a-zA-Z0-9 ]*$"
	_, err := regexp.MatchString(pattern, val)

	if err != nil {
		return false
	}

	if len(val) > 25 {
		return false
	}

	return true
}

func StorePublicData(store *models.Store) *gin.H {
	return &gin.H{
		"name":      store.Name,
		"linkName":  store.LinkName,
		"createdAt": store.CreatedAt,
		"itemCount": store.ItemCount,
		"userId":    store.UserId,
	}
}
