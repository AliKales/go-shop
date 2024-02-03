package tests

import (
	"example/web-service-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

func TestEmailSend(c *gin.Context) {
	utils.SendEmail("ali.kales@hotmail.com", "Deneme", "Test")
}