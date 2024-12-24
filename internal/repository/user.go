package repository

import (
	_ "github.com/lib/pq"
	"gw-currency-wallet/internal/domain"
	"log"
)

func CreateUser(user *domain.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}
