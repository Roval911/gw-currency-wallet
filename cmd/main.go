package main

import (
	_ "gw-currency-wallet/docs"
	"gw-currency-wallet/internal/app"
	"log"
)

// @title Currency Wallet API
// @version 1.0
// @description API for managing wallets and currency exchanges.
// @contact.name API Support
// @contact.email support@currencywallet.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
