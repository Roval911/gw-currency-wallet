package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Хранится в зашифрованном виде
}

type Wallet struct {
	UserID int64   `json:"user_id"`
	USD    float64 `json:"usd"`
	RUB    float64 `json:"rub"`
	EUR    float64 `json:"eur"`
}

type Transaction struct {
	ID       int64   `json:"id"`
	UserID   int64   `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"` // Например, "USD", "RUB", "EUR"
	Type     string  `json:"type"`     // Например, "deposit", "withdraw", "exchange"
	Status   string  `json:"status"`   // Например, "success", "failed"
}

type ExchangeRate struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rate         float64 `json:"rate"`
	UpdatedAt    string  `json:"updated_at"`
}
