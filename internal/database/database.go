package database

import (
	"example/web-service-gin/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB gorm.DB
var Email string
var EmailPassword string

func InitDB() {
	///For local .env use below commented codes!!!!!!!!!

	// wd, errGetwd := os.Getwd()
	// if errGetwd != nil {
	// 	log.Fatal(errGetwd)
	// }

	// wd = strings.Replace(wd, "\\cmd", "", 1)

	// godotenv.Load(wd + "/.env")

	connStr := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = *db

	Email = os.Getenv("EMAIL")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")
}

func GetUserBy(column string, val *string) *models.User {
	if val == nil {
		return nil
	}

	var user models.User
	if err := DB.Where(column+" = ?", val).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func GetStoreBy(column string, val interface{}) *models.Store {
	var store models.Store
	if err := DB.Where(column+" = ?", val).First(&store).Error; err != nil {
		return nil
	}
	return &store
}

func GetStoreItems(storeId int) []models.StoreItem {
	var storeItems []models.StoreItem
	DB.Where("store_id = ?", storeId).Find(&storeItems)
	return storeItems
}

func GetStoreItemBy(column string, val interface{}) *models.StoreItem {
	var item models.StoreItem
	if err := DB.Where(column+" = ?", val).First(&item).Error; err != nil {
		return nil
	}
	return &item
}
