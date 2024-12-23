package db

import (
	"database/sql"
	"log"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CheckUserExists(username, email string) (bool, error) {
	query := `SELECT * FROM users WHERE username = $1 OR email = $2`

	var count int
	err := db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateUser(user *User) error {
	//passwordHash, err := hash.
	//if err != nil {
	//	log.Printf("Ошибка при хэшировании пароля: %v", err)
	//	return fmt.Errorf("could not hash password: %w", err)
	//}

	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}
