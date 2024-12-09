package storage

import (
	"context"

	"github.com/kingxl111/cakes-database-app/internal/models"
)

type Storage struct {
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
	UpdateOrder(userID int, orderID int, paymentMethod string) error
	DeleteOrder(userID, orderID int) error

	GetDeliveryPoints() ([]models.DeliveryPoint, error)
}

type UserCakeManager interface {
	GetCakes() ([]models.Cake, error)
}

type AdminAuthorization interface {
	AddAdmin(username, password string) (int, error)
	GetAdmin(username, password_hash string) (int, error)
}

type Admin interface {
	GetUsers() ([]models.User, error)
	//DeleteUser(userID int) error
	Backup() error
	Restore() error

	AddCake(ctx context.Context, cake models.Cake) (int, error)
	RemoveCake(ctx context.Context, id int) error
}

func NewStorage(db *DB) *Storage {
	return &Storage{
		Authorization:      NewAuthPostgres(db),
		UserOrderManager:   NewUserOrderManagerPostgres(db),
		UserCakeManager:    NewUserCakeManagerPostgres(db),
		AdminAuthorization: NewAdminAuthPostgres(db),
		Admin:              NewAdminPostgres(db),
	}
}
