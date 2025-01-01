package postgres

import (
	"database/sql"
	"errors"
	"fmt"
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

func (s *PostgresStorage) Exchange(userID uint, fromCurrency, toCurrency string, amount, rate float32) (map[string]float32, error) {
	// Начало транзакции
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Проверяем текущий баланс
	var currentBalance float32
	query := `SELECT 
                  CASE 
                    WHEN $2 = 'USD' THEN USD 
                    WHEN $2 = 'RUB' THEN RUB 
                    WHEN $2 = 'EUR' THEN EUR 
                  END as balance 
              FROM wallets WHERE user_id = $1`
	err = tx.QueryRow(query, userID, fromCurrency).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return nil, errors.New("source currency not found")
	} else if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to fetch balance: %w", err)
	}

	if currentBalance < amount {
		tx.Rollback()
		return nil, errors.New("insufficient funds")
	}

	// Обновляем баланс источника (списываем)
	var updateSourceQuery string
	switch fromCurrency {
	case "USD":
		updateSourceQuery = `UPDATE wallets SET USD = USD - $1 WHERE user_id = $2`
	case "RUB":
		updateSourceQuery = `UPDATE wallets SET RUB = RUB - $1 WHERE user_id = $2`
	case "EUR":
		updateSourceQuery = `UPDATE wallets SET EUR = EUR - $1 WHERE user_id = $2`
	default:
		tx.Rollback()
		return nil, fmt.Errorf("unsupported currency: %s", fromCurrency)
	}

	_, err = tx.Exec(updateSourceQuery, amount, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update source balance: %w", err)
	}

	// Рассчитываем сумму для добавления в целевую валюту
	exchangedAmount := amount * rate

	// Обновляем баланс целевой валюты
	var updateTargetQuery string
	switch toCurrency {
	case "USD":
		updateTargetQuery = `UPDATE wallets SET USD = USD + $1 WHERE user_id = $2`
	case "RUB":
		updateTargetQuery = `UPDATE wallets SET RUB = RUB + $1 WHERE user_id = $2`
	case "EUR":
		updateTargetQuery = `UPDATE wallets SET EUR = EUR + $1 WHERE user_id = $2`
	default:
		tx.Rollback()
		return nil, fmt.Errorf("unsupported currency: %s", toCurrency)
	}

	_, err = tx.Exec(updateTargetQuery, exchangedAmount, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update target balance: %w", err)
	}

	// Завершаем транзакцию
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Возвращаем новый баланс пользователя
	balances := make(map[string]float32)
	rows, err := s.db.Query(`SELECT 'USD' as currency, USD as balance FROM wallets WHERE user_id = $1
	                          UNION ALL
	                          SELECT 'RUB', RUB FROM wallets WHERE user_id = $1
	                          UNION ALL
	                          SELECT 'EUR', EUR FROM wallets WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated balances: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var currency string
		var balance float32
		if err := rows.Scan(&currency, &balance); err != nil {
			return nil, fmt.Errorf("failed to scan balances: %w", err)
		}
		balances[currency] = balance
	}

	return balances, nil
}
