package main

import (
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/routes"
	"gw-currency-wallet/internal/storages/postgres"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Printf("Загрузка конфигураций: %+v", cfg)

	db, err := postgres.NewPostgresConnection(postgres.ConnectionInfo{
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

	postgres.SetDB(db)

	postgres.RunMigrations()

	storage := postgres.NewPostgresStorage(db)

	authHandler := hanlers.NewAuthHandler(storage)

	router := routes.SetupRouter(authHandler)

	port := cfg.Server.Port
	log.Printf("Запуск сервера на порту: %d", port)

	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
