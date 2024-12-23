package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gw-currency-wallet/pkg/database"
	"gw-currency-wallet/pkg/migrate"
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
	migrate.RunMigrations()
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(os.Getenv("PORT"))
}
