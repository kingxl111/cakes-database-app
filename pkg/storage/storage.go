package storage

import "cakes-database-app/pkg/models"

type Storage struct {
	Authorization
	OrderManager
}

type Authorization interface {
	CreateUser(user models.User) (int, error)
}

type OrderManager interface {
	CreateOrder(order models.Order) (int, error)
	GetOrder(userID, orderID int) (models.Order, error)
	UpdateOrder(userID int, orderID int, paymentMethod models.Order) error 
	DeleteOrder(userID, orderID int) error
}

func NewStorage(/*db ...*/) *Storage {
	return &Storage{}
}