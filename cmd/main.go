package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/pkg/database"
	"log"
	"os"
)

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Инициализация базы данных
	database.InitDb()

	// Запуск миграций
	database.RunMigrations()
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/api/v1/register", hanlers.CreateUserHandler)

	router.Run(os.Getenv("PORT"))
}
