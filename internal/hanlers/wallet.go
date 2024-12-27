package hanlers

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/storages"
	"net/http"
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
