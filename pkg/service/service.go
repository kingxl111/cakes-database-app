package service

import "cakes-database-app/pkg/models"

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
	CreateOrder(order models.Order) (int, error)
	GetOrder(userID, orderID int) (models.Order, error)
	UpdateOrder(userID int, orderID int, paymentMethod models.Order) error 
	DeleteOrder(userID, orderID int) error
}

