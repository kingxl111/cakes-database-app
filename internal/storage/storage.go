package storage

import (
	"github.com/kingxl111/cakes-database-app/internal/models"
)

type Storage struct {
	Logger

	Authorization
	UserOrderManager
	UserCakeManager

	AdminAuthorization
	Admin
}

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password_hash string) (int, error)
}

type UserOrderManager interface {
	CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error)
	GetOrders(userID int) (models.GetOrdersResponse, error)
	// UpdateOrder(userID int, orderID int, paymentMethod models.Order) error
	// DeleteOrder(userID, orderID int) error
}

type UserCakeManager interface {
	GetCakes() ([]models.Cake, error)
}

type Logger interface {
	WriteLog(level string, msg string) error
}

type AdminAuthorization interface {
	GetAdmin(username, password_hash string) (int, error)
}

type Admin interface {
	GetUsers() ([]models.User, error)
}

func NewStorage(db *DB) *Storage {
	return &Storage{
		Authorization:      NewAuthPostgres(db),
		UserOrderManager:   NewUserOrderManagerPostgres(db),
		UserCakeManager:    NewUserCakeManagerPostgres(db),
		Logger:             NewLoggerPostgres(db),
		AdminAuthorization: NewAdminAuthPostgres(db),
		Admin:              NewAdminPostgres(db),
	}
}
