package main

import (
	"gw-currency-wallet/internal/app"
	"log"
)

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
