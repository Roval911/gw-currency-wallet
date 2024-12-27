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

	public := router.Group("/api/v1")
	{
		public.POST("/register", authHandler.CreateUserHandler)
		public.POST("/login", authHandler.Login)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/createwallet", authHandler.CreateWalletHandle)
		protected.GET("/balance", authHandler.GetBalanceHandle)
		protected.POST("/wallet/deposit", authHandler.DepositHandle)
		protected.POST("/wallet/withdraw", authHandler.WithdrawHandle)
	}

	return router
}
