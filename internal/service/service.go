package service

import (
	"context"

	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
	"github.com/kingxl111/cakes-database-app/internal/storage/s3"
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
	UpdateOrder(userID int, orderID int, paymentMethod string) error
	DeleteOrder(userID, orderID int) error

	GetDeliveryPoints() ([]models.DeliveryPoint, error)
}

type CakeManager interface {
	GetCakes() ([]models.Cake, error)
	GetCake(id int) (models.Cake, error)
}

type AdminAuthorization interface {
	GenerateAdminToken(username, password string) (string, error)
	ParseAdminToken(accessToken string) (int, error)
}

type AdminService interface {
	GetUsers() ([]models.User, error)
	Backup() error
	Restore() error

	AddCake(ctx context.Context, cake models.Cake) (int, error)
	RemoveCake(ctx context.Context, id int) error
}

func NewService(storage *storage.Storage, s s3.ClientS3) *Service {
	return &Service{
		Authorization:      NewAuthService(storage),
		OrderManager:       NewOrderService(storage),
		CakeManager:        NewCakeService(storage, s),
		AdminAuthorization: NewAdminAuthService(storage),
		AdminService:       NewAdminService(storage),
	}
}
