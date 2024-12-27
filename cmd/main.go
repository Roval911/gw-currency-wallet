package main

import (
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/routes"
	"gw-currency-wallet/internal/storages"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Printf("Загрузка конфигураций: %+v", cfg)

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

	storages.SetDB(db)

	storages.RunMigrations()

	storage := storages.NewPostgresStorage(db)

	authHandler := hanlers.NewAuthHandler(storage)

	router := routes.SetupRouter(authHandler)

	port := cfg.Server.Port
	log.Printf("Запуск сервера на порту: %d", port)

	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
