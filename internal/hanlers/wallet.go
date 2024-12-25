package hanlers

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/repository"
	"net/http"
)

func CreateWalletHandle(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := repository.CreateWallet(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать кошелек"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Wallet registered successfully"})
}
