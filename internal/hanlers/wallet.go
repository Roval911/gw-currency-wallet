package hanlers

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/storages"
	"net/http"
)

func (h *AuthHandler) CreateWalletHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.storage.CreateWallet(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать кошелек"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Wallet registered successfully"})
}

func (h *AuthHandler) GetBalanceHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	wallet, err := h.storage.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": wallet})
}

func (h *AuthHandler) DepositHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req storages.DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	wallet, err := h.storage.Deposit(userID, req.Amount, req.Currency)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	wallet, err := h.storage.Withdraw(userID, req.Amount, req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdrawal successful",
		"new_balance": wallet,
	})
}
