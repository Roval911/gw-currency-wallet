package main

import (
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/routes"
	"gw-currency-wallet/internal/storages"
	"log"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded config: %+v", cfg)

	// init db
	db, err := storages.NewPostgresConnection(storages.ConnectionInfo{
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

	// Устанавливаем db в глобальную переменную storages
	storages.SetDB(db)

	// Запускаем миграции
	storages.RunMigrations()

	storage := storages.NewPostgresStorage(db)

	// Создаем обработчик
	authHandler := hanlers.NewAuthHandler(storage)

	router := routes.SetupRouter(authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Дефолтный порт
	}

	router.Run(port)
}
