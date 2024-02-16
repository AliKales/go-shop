package handlers

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/middlewares"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StoreItemReq struct {
	Name      string
	Price     float64
	ItemCount int `json:"itemCount"`
}

func CreateStoreHandler(c *gin.Context) {
	user := middlewares.ConvertToUser(c)

	req := middlewares.ConvertToStoreCreateReq(c)

	linkName := strings.TrimSpace(req.LinkName)

	linkName = strings.ReplaceAll(linkName, " ", "-")

	linkName = strings.ToLower(linkName)

	store := models.Store{UserId: int(user.Id), CreatedAt: time.Now(), Name: req.Name, LinkName: linkName, ItemCount: 0, ItemSellCount: 0}

	if err := database.DB.Create(&store).Error; err != nil {
		if utils.IsNotUniqueColumn(err, "stores_user_id_key") {
			c.JSON(http.StatusConflict, gin.H{"message": "This user already has a store!", "val": "user_id"})
			return
		} else if utils.IsNotUniqueColumn(err, "stores_name_key") {
			c.JSON(http.StatusConflict, gin.H{"message": "Name is already taken", "val": "name"})
			return
		} else if utils.IsNotUniqueColumn(err, "users_link_name_key") {
			c.JSON(http.StatusConflict, gin.H{"message": "Link Name is already taken", "val ": "link_name"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create store"})
		return
	}

	storeId := int(store.ID)
	user.StoreId = &storeId
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "store created!"})
	return
}

func GetStorePublicHandler(c *gin.Context) {
	storeLinkName := c.Param("storeLinkName")

	store := database.GetStoreBy("link_name", storeLinkName)

	if store == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Store not found!"})
		return
	}

	items := database.GetStoreItems(int(store.ID))

	token := middlewares.GetTokenFromHeader(c)
	user := database.GetUserBy("token", token)

	if user == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Store found!", "store": store.PublicData(), "items": items})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store found!", "store": store.PublicData(), "items": items})
}

func CreateStoreItemHandler(c *gin.Context) {
	user := middlewares.ConvertToUser(c)

	if user.StoreId == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have a store!"})
		return
	}

	var req StoreItemReq
	c.BindJSON(&req)

	storeItem := models.StoreItem{StoreId: *user.StoreId, Name: req.Name, CreatedAt: time.Now().UTC(), ItemSellCount: 0, ItemCount: req.ItemCount, Price: req.Price}

	if err := database.DB.Create(&storeItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Store item could not created!"})
		return
	}

	sqlStatement := `
        UPDATE stores
        SET item_count = item_count + 1
        WHERE id = ?
    `

	database.DB.Exec(sqlStatement, uint(*user.StoreId))

	c.JSON(http.StatusOK, gin.H{"message": "Item created!"})
}
