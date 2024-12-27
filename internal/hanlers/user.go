package hanlers

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/middleware"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/hash"
	"net/http"
)

func (h *AuthHandler) CreateUserHandler(c *gin.Context) {
	var user storages.User

	// Чтение JSON из тела запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Printf("Неверные данные ввода при создании пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		h.logger.Printf("Не удалось захешировать пароль для пользователя %s: %v", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = hashedPassword

	if err := h.storage.CreateUser(&user); err != nil {
		h.logger.Printf("Не удалось создать пользователя %s: %v", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest storages.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.logger.Printf("Неверные данные ввода при попытке входа: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Поиск пользователя в базе данных по email
	user, err := h.storage.GetUserByUsername(loginRequest.Username)
	if err != nil {
		h.logger.Printf("Неудачная попытка входа: неверное имя пользователя %s или пароль", loginRequest.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Проверка пароля
	if !hash.CheckPassword(loginRequest.Password, user.Password) {
		h.logger.Printf("Неудачная попытка входа для пользователя %s: неверный пароль", loginRequest.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := middleware.GenerateJWT(user.ID, user.Username)
	if err != nil {
		h.logger.Printf("Не удалось сгенерировать токен для пользователя %s: %v", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

	// Генерация ответа при успешной авторизации
	//c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID})
}
