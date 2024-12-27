package hanlers

import (
	"github.com/sirupsen/logrus"
	"gw-currency-wallet/internal/storages"
)

type AuthHandler struct {
	storage storages.Storages
	logger  *logrus.Logger
}

func NewAuthHandler(storage storages.Storages) *AuthHandler {
	return &AuthHandler{storage: storage}
}
