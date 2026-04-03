package app

import (
	"github.com/S4F4Y4T/goWebService/internal/handler"
)

type App struct {
	UserHandler    *handler.UserHandler
	ProductHandler *handler.ProductHandler
}
