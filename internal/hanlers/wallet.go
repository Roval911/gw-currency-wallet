package hanlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/roval911/proto-exchange/exchange"
	"github.com/sirupsen/logrus"
	"gw-currency-wallet/internal/storages"
	"net/http"
	"time"
)

func (h *AuthHandler) CreateWalletHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.storage.CreateWallet(userID); err != nil {
		h.logger.Printf("Не удалось создать кошелек для пользователя %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать кошелек"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Wallet registered successfully"})
}

func (h *AuthHandler) GetBalanceHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	wallet, err := h.storage.GetBalance(userID)
	if err != nil {
		h.logger.Printf("Не удалось получить баланс для пользователя %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": wallet})
}

func (h *AuthHandler) DepositHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req storages.DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		h.logger.Printf("Неверные данные запроса пополнения для пользователя %d: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	wallet, err := h.storage.Deposit(userID, req.Amount, req.Currency)
	if err != nil {
		h.logger.Printf("Не удалось пополнить кошелек пользователя %d на %f %s: %v", userID, req.Amount, req.Currency, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Account topped up successfully",
		"new_balance": wallet,
	})
}

func (h *AuthHandler) WithdrawHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req storages.WithdrawRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		h.logger.Printf("Неверные данные запроса на снятие для пользователя %d: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	wallet, err := h.storage.Withdraw(userID, req.Amount, req.Currency)
	if err != nil {
		h.logger.Printf("Не удалось снять %f %s с кошелька пользователя %d: %v", req.Amount, req.Currency, userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdrawal successful",
		"new_balance": wallet,
	})
}

type ExchangeHandler struct {
	storage         storages.Storages
	logger          *logrus.Logger
	exchangeService exchange_grpc.ExchangeServiceClient
}

// NewExchangeHandler создает новый обработчик для обмена валют
func NewExchangeHandler(storage storages.Storages, logger *logrus.Logger, exchangeService exchange_grpc.ExchangeServiceClient) *ExchangeHandler {
	return &ExchangeHandler{
		storage:         storage,
		logger:          logger,
		exchangeService: exchangeService,
	}
}

// GetExchangeRatesHandle обрабатывает запрос на получение курсов валют
func (h *ExchangeHandler) GetExchangeRatesHandle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := h.exchangeService.GetExchangeRates(ctx, &exchange_grpc.Empty{})
	if err != nil {
		h.logger.Printf("Не удалось получить курсы валют: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exchange rates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rates": resp.Rates})
}

// ExchangeHandle обрабатывает запрос на обмен валют
func (h *ExchangeHandler) ExchangeHandle(c *gin.Context) {
	var req struct {
		FromCurrency string  `json:"from_currency" binding:"required"`
		ToCurrency   string  `json:"to_currency" binding:"required"`
		Amount       float32 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Printf("Неверный запрос обмена: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверяем, что валюты допустимы
	if !isValidCurrency(req.FromCurrency) || !isValidCurrency(req.ToCurrency) {
		h.logger.Printf("Неверная валюта: %s или %s", req.FromCurrency, req.ToCurrency)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Получаем курс обмена через gRPC
	resp, err := h.exchangeService.GetExchangeRateForCurrency(ctx, &exchange_grpc.CurrencyRequest{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
	})
	if err != nil {
		h.logger.Printf("Не удалось получить курс валют: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exchange rate"})
		return
	}

	// Проверяем баланс и обновляем его
	rate := resp.Rate
	userID := c.GetUint("user_id")
	wallet, err := h.storage.Exchange(userID, req.FromCurrency, req.ToCurrency, req.Amount, rate)
	if err != nil {
		h.logger.Printf("Ошибка обмена валют: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Exchange successful",
		"exchanged_amount": req.Amount * rate,
		"new_balance":      wallet,
	})
}

// Проверка валидности валюты
func isValidCurrency(currency string) bool {
	validCurrencies := []string{"USD", "RUB", "EUR"}
	for _, valid := range validCurrencies {
		if currency == valid {
			return true
		}
	}
	return false
}
