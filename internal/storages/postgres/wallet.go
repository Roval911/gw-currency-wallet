package postgres

import (
	"database/sql"
	"errors"
	"gw-currency-wallet/internal/storages"
)

func (s *PostgresStorage) CreateWallet(userid uint) error {
	query := "INSERT INTO wallets (user_id, USD, RUB, EUR) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(query, userid, 0.0, 0.0, 0.0)
	if err != nil {
		s.logger.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}

func (s *PostgresStorage) GetBalance(userID uint) (storages.Wallet, error) {
	var wallet storages.Wallet
	wallet.UserID = userID // Установить userID вручную, если нужно
	err := db.QueryRow(`
        SELECT USD, RUB, EUR FROM wallets WHERE user_id = $1`, userID).
		Scan(&wallet.USD, &wallet.RUB, &wallet.EUR)
	if err == sql.ErrNoRows {
		s.logger.Printf("Кошелек для пользователя %d не найден", userID)
		return storages.Wallet{}, errors.New("wallet not found")
	}
	if err != nil {
		s.logger.Printf("Ошибка при получении баланса пользователя %d: %v", userID, err)
		return storages.Wallet{}, err
	}
	return wallet, nil
}

// Deposit пополняет баланс пользователя в указанной валюте.
func (s *PostgresStorage) Deposit(userID uint, amount float64, currency string) (storages.Wallet, error) {
	if currency != "USD" && currency != "RUB" && currency != "EUR" {
		s.logger.Printf("Ошибка: неверная валюта при пополнении для пользователя %d: %s", userID, currency)
		return storages.Wallet{}, errors.New("invalid currency")
	}

	query := `UPDATE wallets SET ` + currency + ` = ` + currency + ` + $1 WHERE user_id = $2 RETURNING USD, RUB, EUR`
	var wallet storages.Wallet
	err := db.QueryRow(query, amount, userID).Scan(&wallet.USD, &wallet.RUB, &wallet.EUR)
	if err == sql.ErrNoRows {
		s.logger.Printf("Кошелек для пользователя %d не найден при пополнении на сумму %.2f %s", userID, amount, currency)
		return storages.Wallet{}, errors.New("wallet not found")
	}
	if err != nil {
		s.logger.Printf("Ошибка при пополнении кошелька пользователя %d на сумму %.2f %s: %v", userID, amount, currency, err)
		return storages.Wallet{}, err
	}
	return wallet, nil
}

// Withdraw снимает указанную сумму из баланса пользователя.
func (s *PostgresStorage) Withdraw(userID uint, amount float64, currency string) (storages.Wallet, error) {
	if currency != "USD" && currency != "RUB" && currency != "EUR" {
		s.logger.Printf("Ошибка: неверная валюта при снятии для пользователя %d: %s", userID, currency)
		return storages.Wallet{}, errors.New("invalid currency")
	}

	var currentBalance float64
	err := db.QueryRow(`SELECT `+currency+` FROM wallets WHERE user_id = $1`, userID).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		s.logger.Printf("Кошелек для пользователя %d не найден при снятии суммы %.2f %s", userID, amount, currency)
		return storages.Wallet{}, errors.New("wallet not found")
	}
	if err != nil {
		s.logger.Printf("Ошибка при проверке баланса кошелька пользователя %d для снятия: %v", userID, err)
		return storages.Wallet{}, err
	}

	if currentBalance < amount {
		s.logger.Printf("Недостаточно средств для пользователя %d при попытке снять %.2f %s", userID, amount, currency)
		return storages.Wallet{}, errors.New("insufficient funds")
	}

	query := `UPDATE wallets SET ` + currency + ` = ` + currency + ` - $1 WHERE user_id = $2 RETURNING USD, RUB, EUR`
	var wallet storages.Wallet
	err = db.QueryRow(query, amount, userID).Scan(&wallet.USD, &wallet.RUB, &wallet.EUR)
	if err != nil {
		s.logger.Printf("Ошибка при снятии суммы %.2f %s с кошелька пользователя %d: %v", amount, currency, userID, err)
		return storages.Wallet{}, err
	}
	return wallet, nil
}
