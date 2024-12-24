package hanlers

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/db"
	"gw-currency-wallet/internal/domain"
	"gw-currency-wallet/pkg/hash"
	"net/http"
)

func CreateUserHandler(c *gin.Context) {
	var user domain.SignUp

	// Чтение JSON из тела запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = hashedPassword

	if err := db.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
