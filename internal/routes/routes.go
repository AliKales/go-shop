package routes

import (
	"example/web-service-gin/internal/handlers"
	"example/web-service-gin/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://shopi-rp46.onrender.com"}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// POSTS
	apiPost := r.Group("/api")
	{
		apiPost.POST("/signup", handlers.SignupHandler)
		apiPost.POST("/login", handlers.LoginHandler)
		apiPost.POST("/delete-account", handlers.RequestDeleteAccountHandler)
		apiPost.POST("/create-store", middlewares.CreateStoreMiddleware, middlewares.UserAuthMiddleware, handlers.CreateStoreHandler)
		apiPost.POST("/create-store-item", middlewares.UserAuthMiddleware, handlers.CreateStoreItemHandler)
		apiPost.POST("/change-password", handlers.ChangePasswordHandler)
	}

	// GETS
	apiGet := r.Group("/api")
	{
		apiGet.GET("/refresh-token", handlers.RefreshTokenHandler)
		apiGet.GET("/user-action", handlers.RedirectUserAction)
		apiGet.GET("/forget-password", handlers.ForgetPasswordHandler)
		apiGet.GET("/user/:username", handlers.GetUserPublicHandler)
		apiGet.GET("/store/:storeLinkName", handlers.GetStorePublicHandler)
		apiGet.GET("/cart", middlewares.UserAuthMiddleware, handlers.GetCartHandler)
		apiGet.GET("/add-to-cart", middlewares.UserAuthMiddleware, handlers.AddItemToCartHandler)
		apiGet.GET("/delete-from-cart", middlewares.UserAuthMiddleware, handlers.DeleteItemFromCartHandler)
	}

	return r
}
