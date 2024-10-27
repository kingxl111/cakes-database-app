package storage

import (
	"cakes-database-app/pkg/models"
	"context"
	"fmt"
	"time"
)

type OrderPostgres struct {
	db *DB
}

func NewOrderPostgres(db *DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (o *OrderPostgres) CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error) {
	var orderID int
	// tx begin
	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %v", err)
	}
	// if return will be successful, this defer never call
	defer func() {
        if err != nil {
            tx.Rollback(ctx) 
        }
    }()

	// adding new order
	order := models.Order{
        Time:          	time.Now().Format(time.RFC3339),
        UserID:       	userID,
        OrderStatus:   	"pending", 
        PaymentMethod: 	paymentMethod,
    }
	orderQuery := `INSERT INTO orders (time, order_status, user_id,  payment_method) VALUES ($1, $2, $3, $4) RETURNING id`
    err = tx.QueryRow(ctx, 
		orderQuery, 
		order.Time, 
		order.OrderStatus, 
		order.UserID, 
		order.PaymentMethod).Scan(&orderID)
    if err != nil {
        return 0, fmt.Errorf("could not insert order: %v", err)
    }

	var deliveryID int // the same with orderID
    deliveryQuery := `INSERT INTO deliveries (point_id, cost, status, weight) VALUES ($1, $2, $3, $4) RETURNING id`
    err = tx.QueryRow(ctx, 
		deliveryQuery, 
		delivery.PointID, 
		delivery.Cost, 
		delivery.Status, 
		delivery.Weight).Scan(&deliveryID)
    if err != nil {
        return 0, fmt.Errorf("could not insert delivery: %v", err)
    }

	for _, cake := range cakes {
        orderCake := models.OrderCake{
            OrderID: orderID,
            CakeID:  cake.ID,
        }
        orderCakeQuery := `INSERT INTO order_cakes (order_id, cake_id) VALUES ($1, $2)`
        _, err = tx.Exec(ctx, 
			orderCakeQuery, 
			orderCake.OrderID, 
			orderCake.CakeID)
        if err != nil {
            return 0, fmt.Errorf("could not insert order-cake relation: %v", err)
        }
    }

    err = tx.Commit(ctx)
	if err != nil {
        return 0, fmt.Errorf("could not commit transaction: %v", err)
    }

	return orderID, nil
}