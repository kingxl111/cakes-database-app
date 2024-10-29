package service

import (
	"cakes-database-app/pkg/models"
	"cakes-database-app/pkg/storage"

)

type Service struct {
	Authorization
	OrderManager
}

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type OrderManager interface {
	CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error)
	GetOrders(userID int) (models.GetOrdersResponse, error)
	// UpdateOrder(userID int, orderID int, paymentMethod models.Order) error 
	// DeleteOrder(userID, orderID int) error
}

func NewService(storage *storage.Storage) *Service{
	return &Service{
		Authorization: NewAuthService(storage),
		OrderManager: NewOrderService(storage),
	}
}
