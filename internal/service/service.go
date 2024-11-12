package service

import (
	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
)

type Service struct {
	Authorization
	OrderManager
	CakeManager

	AdminAuthorization
	AdminService
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
	DeleteOrder(userID, orderID int) error
}

type CakeManager interface {
	GetCakes() ([]models.Cake, error)
}

type AdminAuthorization interface {
	GenerateAdminToken(username, password string) (string, error)
	ParseAdminToken(accessToken string) (int, error)
}

type AdminService interface {
	GetUsers() ([]models.User, error)
	Backup() error
	Restore() error
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Authorization:      NewAuthService(storage),
		OrderManager:       NewOrderService(storage),
		CakeManager:        NewCakeService(storage),
		AdminAuthorization: NewAdminAuthService(storage),
		AdminService:       NewAdminService(storage),
	}
}
