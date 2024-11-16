package service

import (
	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
)

type OrderService struct {
	stg storage.UserOrderManager
}

func NewOrderService(stg storage.UserOrderManager) *OrderService {
	return &OrderService{stg: stg}
}

func (o *OrderService) CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error) {
	return o.stg.CreateOrder(userID, delivery, cakes, paymentMethod)
}

func (o *OrderService) GetOrders(userID int) (models.GetOrdersResponse, error) {
	return o.stg.GetOrders(userID)
}

func (o *OrderService) DeleteOrder(userID, orderID int) error {
	return o.stg.DeleteOrder(userID, orderID)
}

func (o *OrderService) GetDeliveryPoints() ([]models.DeliveryPoint, error) {
	return o.stg.GetDeliveryPoints()
}

func (o *OrderService) UpdateOrder(userID int, orderID int, paymentMethod string) error {
	return o.stg.UpdateOrder(userID, orderID, paymentMethod)
}
