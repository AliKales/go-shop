package handlers

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/middlewares"
	"example/web-service-gin/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RedirectUserAction(c *gin.Context) {
	tokenParam := c.Query("token")

	user := database.GetUserBy("user_action", &tokenParam)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Token does not exist!"})
		return
	}

	if utils.IsExpired(user.UserActionExpireAt) {
		c.JSON(http.StatusGone, gin.H{"message": "Token expired!"})
		return
	}

	action := strings.Split(user.UserAction, "|")[1]

	if action == "verify-email" {
		VerifyEmailHandler(c, *user)
		return
	} else if action == "delete-account" {
		DeleteAccountHandler(c, *user)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unexpected action!"})
		return
	}
}

func GetUserPublicHandler(c *gin.Context) {
	username := c.Param("username")

	user := database.GetUserBy("username", &username)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	store := database.GetStoreBy("user_id", int(user.Id))

	if store != nil {
		c.JSON(http.StatusOK, gin.H{"message": "User found!", "user": user.PublicData(), "store": store.PublicData()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User found!", "user": user.PublicData()})
}

func AddItemToCartHandler(c *gin.Context) {
	user := middlewares.ConvertToUser(c)

	itemId := c.Query("item")
	i, _ := strconv.Atoi(itemId)

	item := database.GetStoreItemBy("id", uint(i))

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found!"})
		return
	}

	sqlStatement := `
	UPDATE users
	SET cart = jsonb_set(
		cart,
		?,
		(CASE 
			WHEN cart->? IS NOT NULL THEN 
				to_jsonb((cart->>?)::int + 1)::text::jsonb 
			ELSE 
				'1'::jsonb 
		END)
	)
	WHERE id = ?;
	`
	itemIdString := strconv.FormatUint(uint64(item.ID), 10)

	database.DB.Exec(sqlStatement, "{"+itemIdString+"}", itemIdString, itemIdString, user.Id)

	c.JSON(http.StatusOK, gin.H{"message": "Item added to your cart!"})
}

func DecreaseItemCountFromCartHandler(c *gin.Context) {
	user := middlewares.ConvertToUser(c)

	itemId := c.Query("item")

	if !utils.IsOnlyNumbers(itemId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only numbers allowed!"})
		return
	}

	sqlStatement := `
	UPDATE users
	SET cart = 
		CASE 
			WHEN cart->? IS NOT NULL AND (cart->>?)::int > 1 THEN 
				jsonb_set(cart, ?, to_jsonb((cart->>?)::int - 1)::text::jsonb)
			WHEN cart->? IS NOT NULL AND (cart->>?)::int = 1 THEN 
				cart - ?
			ELSE 
				cart
		END
	WHERE id = ?;
	`
	database.DB.Exec(sqlStatement, itemId, itemId, "{"+itemId+"}", itemId, itemId, itemId, itemId, user.Id)

	c.JSON(http.StatusOK, gin.H{"message": "Item decreased!"})
}

func DeleteItemFromCartHandler(c *gin.Context) {
	user := middlewares.ConvertToUser(c)

	itemId := c.Query("item")

	if !utils.IsOnlyNumbers(itemId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only numbers allowed!"})
		return
	}

	sqlStatement := `
	UPDATE users
	SET cart = 
		CASE 
			WHEN cart->? IS NOT NULL THEN 
				cart - ?
			ELSE 
				cart
		END
	WHERE id = ?;
	`
	database.DB.Exec(sqlStatement, itemId, itemId, user.Id)

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted from your cart!"})
}

func GetCartHandler(c *gin.Context) {
	user := middlewares.ConvertToUser(c)

	var cart map[string]int
	user.Cart.AssignTo(&cart)

	resultList := []map[string]interface{}{}

	totalPrice := 0.0

	for key, value := range cart {
		i, _ := strconv.Atoi(key)
		item := database.GetStoreItemBy("id", uint(i))

		if item != nil {
			newObject := map[string]interface{}{"count": value, "price": item.Price, "id": key, "name": item.Name}
			resultList = append(resultList, newObject)
			totalPrice += item.Price
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resultList, "totalPrice": totalPrice, "itemLength": len(resultList)})
}
