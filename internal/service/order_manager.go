package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
	"time"
)

const ordersChan = "orders:events"

type OrderService struct {
	stg   storage.UserOrderManager
	redis *redis.Client
}

func NewOrderService(stg storage.UserOrderManager, rdb *redis.Client) *OrderService {
	return &OrderService{
		stg:   stg,
		redis: rdb,
	}
}

func (o *OrderService) CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error) {
	orderID, err := o.stg.CreateOrder(userID, delivery, cakes, paymentMethod)
	if err != nil {
		return 0, err
	}

	evt := struct {
		Event     string `json:"event"`
		OrderID   int    `json:"order_id"`
		UserID    int    `json:"user_id"`
		Timestamp int64  `json:"ts"`
	}{
		Event:     "new_order",
		OrderID:   orderID,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
	}
	if payload, err := json.Marshal(evt); err == nil {
		o.redis.Publish(context.Background(), ordersChan, payload)
	} else {
		fmt.Printf("pubsub: marshal new_order failed: %v\n", err)
	}

	return orderID, nil
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
	if err := o.stg.UpdateOrder(userID, orderID, paymentMethod); err != nil {
		return err
	}

	evt := struct {
		Event      string `json:"event"`
		OrderID    int    `json:"order_id"`
		UserID     int    `json:"user_id"`
		NewPayment string `json:"payment_method"`
		Timestamp  int64  `json:"ts"`
	}{
		Event:      "order_updated",
		OrderID:    orderID,
		UserID:     userID,
		NewPayment: paymentMethod,
		Timestamp:  time.Now().Unix(),
	}
	if payload, err := json.Marshal(evt); err == nil {
		o.redis.Publish(context.Background(), ordersChan, payload)
	} else {
		fmt.Printf("pubsub: marshal order_updated failed: %v\n", err)
	}

	return nil
}
