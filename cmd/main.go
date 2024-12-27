package main

import (
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/routes"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/pkg/logger"
	"log"
)

func main() {
	logger := logger.InitLogger()
	cfg, err := config.New()
	if err != nil {
		logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	logger.Printf("Загрузка конфигураций: %+v", cfg)

	db, err := postgres.NewPostgresConnection(postgres.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	postgres.SetDB(db)

	postgres.RunMigrations()

	storage := postgres.NewPostgresStorage(db)

	authHandler := hanlers.NewAuthHandler(storage)

	router := routes.SetupRouter(authHandler)

	port := cfg.Server.Port
	log.Printf("Запуск сервера на порту: %d", port)

	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
