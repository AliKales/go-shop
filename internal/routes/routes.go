package routes

import (
	"example/web-service-gin/internal/handlers"
	"example/web-service-gin/internal/middlewares"
	"example/web-service-gin/internal/tests"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// POSTS
	r.POST("/api/signup", handlers.SignupHandler)
	r.POST("/api/login", handlers.LoginHandler)
	r.POST("/api/delete-account", handlers.RequestDeleteAccountHandler)
	r.POST("/api/create-store", middlewares.CreateStoreMiddleware, middlewares.UserAuthMiddleware, handlers.CreateStoreHandler)
	r.POST("/api/create-store-item", middlewares.UserAuthMiddleware, handlers.CreateStoreItemHandler)
	r.POST("/api/change-password", handlers.ChangePasswordHandler)

	// GETS
	r.GET("/api/refresh-token", handlers.RefreshTokenHandler)
	r.GET("/api/user-action", handlers.RedirectUserAction)
	r.GET("/api/forget-password", handlers.ForgetPasswordHandler)
	r.GET("/api/user/:username", handlers.GetUserPublicHandler)
	r.GET("/api/store/:storeLinkName", handlers.GetStorePublicHandler)
	r.GET("/api/cart", middlewares.UserAuthMiddleware, handlers.GetCartHandler)
	r.GET("/api/add-to-cart", middlewares.UserAuthMiddleware, handlers.AddItemToCartHandler)
	r.GET("/api/delete-from-cart", middlewares.UserAuthMiddleware, handlers.DeleteItemFromCartHandler)

	//Tests
	r.POST("/api/sendEmail", tests.TestEmailSend)

	return r
}
