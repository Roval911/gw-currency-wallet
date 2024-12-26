package storages

import (
	"log"
)

func (s *PostgresStorage) CreateWallet(userid uint) error {
	query := "INSERT INTO wallets (user_id, USD, RUB, EUR) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(query, userid, 0.0, 0.0, 0.0)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}
