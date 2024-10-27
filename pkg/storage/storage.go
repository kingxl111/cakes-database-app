package storage

import (
	"cakes-database-app/pkg/models"

	// "github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Authorization
	OrderManager
}

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password_hash string) (int, error)
}

type OrderManager interface {
	CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error)
// 	GetOrder(userID, orderID int) (models.Order, error)
// 	UpdateOrder(userID int, orderID int, paymentMethod models.Order) error 
// 	DeleteOrder(userID, orderID int) error
}

func NewStorage(db *DB) *Storage {
	return &Storage{
		Authorization: NewAuthPostgres(db),
		OrderManager: NewOrderPostgres(db),
	}
}