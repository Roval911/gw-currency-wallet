package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/middleware"
	"gw-currency-wallet/internal/repository"
	"log"
	"os"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %+v\n", cfg)

	// init db
	db, err := repository.NewPostgresConnection(repository.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Устанавливаем db в глобальную переменную repository
	repository.SetDB(db)

	// Запускаем миграции
	repository.RunMigrations()

	// Инициализируем сервер
	router := gin.Default()
	router.Use(cors.Default())

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/api/v1/register", hanlers.CreateUserHandler)
	router.POST("/api/v1/login", hanlers.Login)

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/createwallet", hanlers.CreateWalletHandle)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Дефолтный порт
	}

	router.Run(port)
}
