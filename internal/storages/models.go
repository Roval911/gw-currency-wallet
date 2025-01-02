package storages

// swagger:model
type User struct {
	ID       uint   `json:"-"`
	Username string `json:"username" binding:"required,min=3,gte=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6"`
}

// swagger:model
type LoginRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"password123"`
}

// swagger:model
type Wallet struct {
	UserID uint    `json:"user_id"`
	USD    float64 `json:"USD"`
	RUB    float64 `json:"RUB"`
	EUR    float64 `json:"EUR"`
}

// DepositRequest struct for deposit
// swagger:model
type DepositRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,oneof=USD RUB EUR"`
}

// WithdrawRequest struct for withdrawal
// swagger:model
type WithdrawRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,oneof=USD RUB EUR"`
}

// ExchangeRequest struct for exchange request
// swagger:model
type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Amount       float32 `json:"amount" binding:"required,gt=0"`
}
