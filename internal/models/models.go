package models

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required,min=3,gte=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6"`
}

//type Wallet struct {
//	ID     uint    `json:"id"`
//	UserID uint    `json:"user_id"`
//	USD    float64 `json:"USD"`
//	RUB    float64 `json:"RUB"`
//	EUR    float64 `json:"EUR"`
//}

type LoginRequest struct {
	//ID       uint   `json:"id"`
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"password123"`
}

type Wallet struct {
	UserID uint    `json:"user_id"`
	USD    float64 `json:"USD"`
	RUB    float64 `json:"RUB"`
	EUR    float64 `json:"EUR"`
}

type DepositRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type WithdrawRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
