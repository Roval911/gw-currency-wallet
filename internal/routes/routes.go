package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/middleware"
)

func SetupRouter(authHandler *hanlers.AuthHandler) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// Публичные маршруты
	router.POST("/api/v1/register", authHandler.CreateUserHandler)
	router.POST("/api/v1/login", authHandler.Login)

	// Защищенные маршруты
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/createwallet", authHandler.CreateWalletHandle)
	}

	return router
}
