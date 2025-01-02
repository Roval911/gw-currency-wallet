package hanlers

import (
	"github.com/go-redis/redis/v8"
	exchange_grpc "github.com/roval911/proto-exchange/exchange"
	"github.com/sirupsen/logrus"
	"gw-currency-wallet/internal/storages"
)

type AuthHandler struct {
	storage storages.Storages
	logger  *logrus.Logger
}

func NewAuthHandler(storage storages.Storages, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		storage: storage,
		logger:  logger,
	}
}

type ExchangeHandler struct {
	storage         storages.Storages
	logger          *logrus.Logger
	exchangeService exchange_grpc.ExchangeServiceClient
	redisClient     *redis.Client
}

// NewExchangeHandler создает новый обработчик для обмена валют
func NewExchangeHandler(storage storages.Storages, logger *logrus.Logger, exchangeService exchange_grpc.ExchangeServiceClient, redisClient *redis.Client) *ExchangeHandler {
	return &ExchangeHandler{
		storage:         storage,
		logger:          logger,
		exchangeService: exchangeService,
		redisClient:     redisClient,
	}
}
