package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gw-currency-wallet/internal/storages"
	"log"
)

func (s *PostgresStorage) CreateUser(user *storages.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}

func (s *PostgresStorage) GetUserByUsername(email string) (*storages.User, error) {
	user := &storages.User{}
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		return nil, err
	}
	return user, err
}
