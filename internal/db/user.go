package db

import (
	"database/sql"
	"gw-currency-wallet/internal/domain"
	"log"
)

var db *sql.DB

func CreateUser(user *domain.SignUp) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}
