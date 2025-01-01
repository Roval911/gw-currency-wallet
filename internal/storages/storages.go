package storages

type Storages interface {
	//	user
	CreateUser(user *User) error
	GetUserByUsername(email string) (*User, error)

	//	wallet
	CreateWallet(userid uint) error
	GetBalance(userID uint) (Wallet, error)
	Deposit(userID uint, amount float64, currency string) (Wallet, error)
	Withdraw(userID uint, amount float64, currency string) (Wallet, error)
	Exchange(userID uint, fromCurrency, toCurrency string, amount float32, rate float32) (map[string]float32, error)
}
