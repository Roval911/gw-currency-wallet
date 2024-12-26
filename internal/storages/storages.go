package storages

type Storages interface {
	//	user
	CreateUser(user *User) error
	GetUserByUsername(email string) (*User, error)

	//	wallet
	CreateWallet(userid uint) error
}
