package storage

import (
	"cakes-database-app/pkg/models"
)

type Storage struct {
	Logger
	Authorization
	OrderManager
	CakeManager
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

type CakeManager interface {
	GetCakes() ([]models.Cake, error)
}

type Logger interface {
	WriteLog(level string, msg string) error 
}

func NewStorage(db *DB) *Storage {
	return &Storage{
		Authorization: NewAuthPostgres(db),
		OrderManager: NewOrderPostgres(db),
		CakeManager: NewCakeManagerPostgres(db),
		Logger: NewLoggerPostgres(db),
	}
}
