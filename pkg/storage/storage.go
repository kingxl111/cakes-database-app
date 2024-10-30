package storage

import (
	"cakes-database-app/pkg/models"
	"context"
)

type Storage struct {
	Logger
	Authorization
	OrderManager
}

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password_hash string) (int, error)
}

type OrderManager interface {
	CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error)
	GetOrders(userID int) (models.GetOrdersResponse, error)
// 	UpdateOrder(userID int, orderID int, paymentMethod models.Order) error 
// 	DeleteOrder(userID, orderID int) error
}

type Logger interface {
	WriteLog(ctx *context.Context, level string, msg string) error 
}

func NewStorage(db *DB) *Storage {
	return &Storage{
		Authorization: NewAuthPostgres(db),
		OrderManager: NewOrderPostgres(db),
		Logger: NewLoggerPostgres(db),
	}
}
