package middlewares

import (
	"example/web-service-gin/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StoreCreateReq struct {
	Name     string
	LinkName string `json:"linkName"`
}

func CreateStoreMiddleware(c *gin.Context) {
	var req StoreCreateReq
	c.BindJSON(&req)

	if !controllers.IsLinkNameOK(req.LinkName) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Link Name is not suitable", "val": "link_name"})
		c.Abort()
		return
	}

	if len(req.Name) > 20 || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is not suitable", "val": "name"})
		c.Abort()
		return
	}

	c.Set("requestBody", req)

	c.Next()
}

func ConvertToStoreCreateReq(c *gin.Context) *StoreCreateReq {
	user, exists := c.Get("requestBody")
	if !exists {
		return nil
	}

	req, ok := user.(StoreCreateReq)
	if !ok {
		return nil
	}

	return &req
}
