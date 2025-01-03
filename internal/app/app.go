package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/roval911/proto-exchange/exchange"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/hanlers"
	"gw-currency-wallet/internal/routes"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/pkg/logger"
)

// App структура для представления всего приложения
type App struct {
	logger      *logrus.Logger
	router      *gin.Engine
	redisClient *redis.Client
}

// New создаёт новый экземпляр приложения
func New() (*App, error) {
	// Инициализация логгера
	log := logger.InitLogger()

	// Загрузка конфигурации
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
		return nil, err
	}

	// Инициализация подключения к базе данных
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
		return nil, err
	}

	// Настройка базы данных
	postgres.SetDB(db)
	postgres.RunMigrations()

	// Создание хранилища
	storage := postgres.NewPostgresStorage(db)

	conn, err := grpc.Dial(cfg.ExchangeService.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC-сервису: %v", err)
		return nil, err
	}
	exchangeService := exchange_grpc.NewExchangeServiceClient(conn)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,  // Адрес Redis-сервера
		Password: cfg.Redis.Password, // Пароль Redis (если есть)
		DB:       cfg.Redis.DB,       // Номер базы данных
	})
	// Проверяем подключение к Redis
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
		return nil, err
	}

	// Создание обработчиков
	authHandler := hanlers.NewAuthHandler(storage, log)
	exchangeHandler := hanlers.NewExchangeHandler(storage, log, exchangeService, redisClient)

	// Настройка маршрутов
	router := routes.SetupRouter(authHandler, exchangeHandler)

	// Возвращаем структуру приложения с логгером и маршрутизатором
	return &App{
		logger: log,
		router: router,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	cfg, err := config.New()
	if err != nil {
		a.logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
		return err
	}

	port := cfg.Server.Port
	a.logger.Printf("Запуск сервера на порту: %d", port)

	// Запускаем сервер
	err = a.router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		a.logger.Fatalf("Ошибка запуска сервера: %v", err)
		return err
	}

	return nil
}
