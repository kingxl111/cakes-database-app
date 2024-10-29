package service

import (
	"cakes-database-app/pkg/models"
	"cakes-database-app/pkg/storage"
)

type Order struct {
	stg storage.OrderManager
}

func NewOrderService(stg storage.OrderManager) *Order {
	return &Order{stg: stg}
}

func (o *Order) CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error) {
	return o.stg.CreateOrder(userID, delivery, cakes,paymentMethod)
}

func (o *Order) GetOrders(userID int) (models.GetOrdersResponse, error) {
	return o.stg.GetOrders(userID)	
}