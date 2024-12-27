package storages

import (
	"database/sql"
	"errors"
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

func (s *PostgresStorage) GetBalance(userID uint) (Wallet, error) {
	var wallet Wallet
	wallet.UserID = userID // Установить userID вручную, если нужно
	err := db.QueryRow(`
        SELECT USD, RUB, EUR FROM wallets WHERE user_id = $1`, userID).
		Scan(&wallet.USD, &wallet.RUB, &wallet.EUR)
	if err == sql.ErrNoRows {
		return Wallet{}, errors.New("wallet not found")
	}
	if err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

// Deposit пополняет баланс пользователя в указанной валюте.
func (s *PostgresStorage) Deposit(userID uint, amount float64, currency string) (Wallet, error) {
	if currency != "USD" && currency != "RUB" && currency != "EUR" {
		return Wallet{}, errors.New("invalid currency")
	}

	query := `UPDATE wallets SET ` + currency + ` = ` + currency + ` + $1 WHERE user_id = $2 RETURNING USD, RUB, EUR`
	var wallet Wallet
	err := db.QueryRow(query, amount, userID).Scan(&wallet.USD, &wallet.RUB, &wallet.EUR)
	if err == sql.ErrNoRows {
		return Wallet{}, errors.New("wallet not found")
	}
	if err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

// Withdraw снимает указанную сумму из баланса пользователя.
func (s *PostgresStorage) Withdraw(userID uint, amount float64, currency string) (Wallet, error) {
	if currency != "USD" && currency != "RUB" && currency != "EUR" {
		return Wallet{}, errors.New("invalid currency")
	}

	var currentBalance float64
	err := db.QueryRow(`SELECT `+currency+` FROM wallets WHERE user_id = $1`, userID).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		return Wallet{}, errors.New("wallet not found")
	}
	if err != nil {
		return Wallet{}, err
	}

	if currentBalance < amount {
		return Wallet{}, errors.New("insufficient funds")
	}

	query := `UPDATE wallets SET ` + currency + ` = ` + currency + ` - $1 WHERE user_id = $2 RETURNING USD, RUB, EUR`
	var wallet Wallet
	err = db.QueryRow(query, amount, userID).Scan(&wallet.USD, &wallet.RUB, &wallet.EUR)
	if err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}
