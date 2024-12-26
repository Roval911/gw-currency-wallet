package hanlers

import "gw-currency-wallet/internal/storages"

type AuthHandler struct {
	storage storages.Storages
}

func NewAuthHandler(storage storages.Storages) *AuthHandler {
	return &AuthHandler{storage: storage}
}
