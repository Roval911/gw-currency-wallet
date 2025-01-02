package hanlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/roval911/proto-exchange/exchange"
	"gw-currency-wallet/internal/storages"
	"net/http"
	"time"
)

// CreateWalletHandle godoc
// @Summary Create a new wallet
// @Description Creates a wallet for the authenticated user
// @Tags Wallet
// @Produce json
// @Success 201 {object} map[string]interface{} "Wallet registered successfully"
// @Failure 500 {object} map[string]interface{} "Failed to create wallet"
// @Security BearerToken
// @Router /api/v1/createwallet [post]
func (h *AuthHandler) CreateWalletHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.storage.CreateWallet(userID); err != nil {
		h.logger.Printf("Не удалось создать кошелек для пользователя %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать кошелек"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Wallet registered successfully"})
}

// GetBalanceHandle godoc
// @Summary Get wallet balance
// @Description Retrieves the balance of the authenticated user's wallet
// @Tags Wallet
// @Produce json
// @Success 200 {object} storages.Wallet "Balance retrieved"
// @Failure 500 {object} map[string]interface{} "Failed to fetch balance"
// @Security BearerToken
// @Router /api/v1/balance [get]
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

// DepositHandle godoc
// @Summary Deposit funds to wallet
// @Description Adds funds to the authenticated user's wallet
// @Tags Wallet
// @Accept  json
// @Produce  json
// @Param depositRequest body storages.DepositRequest true "Deposit details"
// @Success 200 {object} map[string]interface{} "Account topped up successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to deposit funds"
// @Security BearerToken
// @Router /api/v1/wallet/deposit [post]
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

// WithdrawHandle godoc
// @Summary Withdraw funds from wallet
// @Description Withdraws funds from the authenticated user's wallet
// @Tags Wallet
// @Accept  json
// @Produce  json
// @Param withdrawRequest body storages.WithdrawRequest true "Withdrawal details"
// @Success 200 {object} map[string]interface{} "Withdrawal successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to withdraw funds"
// @Security BearerToken
// @Router /api/v1/wallet/withdraw [post]
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

// GetExchangeRatesHandle godoc
// @Summary Get exchange rates
// @Description Retrieves current exchange rates for supported currencies from the exchange service or cache
// @Tags Exchange
// @Produce json
// @Success 200 {object} map[string]interface{} "Exchange rates retrieved"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve exchange rates"
// @Security BearerToken
// @Router /api/v1/exchange/rates [get]
func (h *ExchangeHandler) GetExchangeRatesHandle(c *gin.Context) {
	// Проверка кеша Redis
	cacheKey := "exchange_rates"
	rateCache, err := h.redisClient.Get(context.Background(), cacheKey).Result()

	if err == nil {
		// Если кэширование успешно (данные найдены в Redis)
		var rates map[string]float64
		if err := json.Unmarshal([]byte(rateCache), &rates); err == nil {
			c.JSON(http.StatusOK, gin.H{"rates": rates})
			return
		}
	}

	// Если нет в кеше или произошла ошибка, запрашиваем у внешнего сервиса
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := h.exchangeService.GetExchangeRates(ctx, &exchange_grpc.Empty{})
	if err != nil {
		h.logger.Printf("Не удалось получить курсы валют: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exchange rates"})
		return
	}

	// Кэшируем данные в Redis
	rateData, err := json.Marshal(resp.Rates)
	if err != nil {
		h.logger.Printf("Не удалось сериализовать курсы валют для кеша: %v", err)
	}

	// Сохраняем в Redis с временем жизни 10 минут
	err = h.redisClient.Set(context.Background(), cacheKey, rateData, time.Minute*10).Err()
	if err != nil {
		h.logger.Printf("Не удалось сохранить курсы валют в Redis: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"rates": resp.Rates})
}

// ExchangeHandle godoc
// @Summary Exchange currency
// @Description Exchanges one currency to another in the authenticated user's wallet
// @Tags Exchange
// @Accept  json
// @Produce  json
// @Param exchangeRequest body storages.ExchangeRequest true "Exchange details"
// @Success 200 {object} map[string]interface{} "Exchange successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to process exchange"
// @Security BearerToken
// @Router /api/v1/exchange [post]
func (h *ExchangeHandler) ExchangeHandle(c *gin.Context) {
	var req storages.ExchangeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Printf("Неверный запрос обмена: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !isValidCurrency(req.FromCurrency) || !isValidCurrency(req.ToCurrency) {
		h.logger.Printf("Неверная валюта: %s или %s", req.FromCurrency, req.ToCurrency)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency"})
		return
	}

	// Проверка кеша Redis для обменного курса
	cacheKey := "exchange_rate_" + req.FromCurrency + "_" + req.ToCurrency
	rateCache, err := h.redisClient.Get(context.Background(), cacheKey).Result()

	var rate float32 // Заменили на float32, чтобы соответствовать типу в структуре

	if err == nil {
		// Если кэширование успешно
		if err := json.Unmarshal([]byte(rateCache), &rate); err == nil {
			// Применяем курс из кеша
		}
	} else {
		// Если курса нет в кэше, делаем запрос к внешнему сервису
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := h.exchangeService.GetExchangeRateForCurrency(ctx, &exchange_grpc.CurrencyRequest{
			FromCurrency: req.FromCurrency,
			ToCurrency:   req.ToCurrency,
		})
		if err != nil {
			h.logger.Printf("Не удалось получить курс валют: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exchange rate"})
			return
		}

		rate = float32(resp.Rate) // Приводим к float32, так как структура ожидает float32

		// Кэшируем курс в Redis
		rateData, err := json.Marshal(rate)
		if err != nil {
			h.logger.Printf("Не удалось сериализовать курс валют для кеша: %v", err)
		}

		// Сохраняем курс в Redis с временем жизни 10 минут
		err = h.redisClient.Set(context.Background(), cacheKey, rateData, time.Minute*10).Err()
		if err != nil {
			h.logger.Printf("Не удалось сохранить курс валют в Redis: %v", err)
		}
	}

	// Производим обмен валют
	userID := c.GetUint("user_id")
	exchangedAmount := float32(req.Amount) * rate // Приводим к float32

	wallet, err := h.storage.Exchange(userID, req.FromCurrency, req.ToCurrency, exchangedAmount, rate)
	if err != nil {
		h.logger.Printf("Ошибка обмена валют: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Exchange successful",
		"exchanged_amount": exchangedAmount,
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
