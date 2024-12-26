package hanlers

import (
	"github.com/gin-gonic/gin"
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
