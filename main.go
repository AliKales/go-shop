package main

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/routes"
	"example/web-service-gin/internal/utils"
)

func init() {
	utils.InitGeo()
	database.InitDB()
}

func main() {
	r := routes.SetupRouter()

	if err := r.Run(":" + database.Port); err != nil {
		panic("failed to start the server")
	}
}