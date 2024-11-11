package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/kingxl111/cakes-database-app/internal/models"
)

type UserOrderManagerPostgres struct {
	db *DB
}

func NewUserOrderManagerPostgres(db *DB) *UserOrderManagerPostgres {
	return &UserOrderManagerPostgres{db: db}
}

func (o *UserOrderManagerPostgres) CreateOrder(userID int, delivery models.Delivery, cakes []models.Cake, paymentMethod string) (int, error) {
	var orderID int
	// tx begin
	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %v", err)
	}

	// calculate common order's cakes cost by cakes ids:
	cost := 0
	for _, cake := range cakes {
		var price int
		id := cake.ID

		builderSelect := sq.Select("price").
			From(CakesTable).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{"id": id})

		query, args, err := builderSelect.ToSql()
		if err != nil {
			return orderID, fmt.Errorf("could not build query: %v", err.Error())
		}

		err = tx.QueryRow(ctx, query, args...).Scan(&price)
		if err != nil {
			return 0, err
		}

		cost += price
	}

	// adding new order
	order := models.Order{
		Time:          time.Now(),
		UserID:        userID,
		OrderStatus:   "pending",
		PaymentMethod: paymentMethod,
		Cost:          cost, // sum of order's cakes cost
	}

	builderInsert := sq.Insert(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Columns("time",
			"order_status",
			"user_id",
			"payment_method",
			"cost").
		Values(order.Time,
			order.OrderStatus,
			order.UserID,
			order.PaymentMethod,
			order.Cost).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		return orderID, fmt.Errorf("could not build query: %v", err.Error())
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&orderID)
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

func (o *UserOrderManagerPostgres) GetOrders(userID int) (models.GetOrdersResponse, error) {
	const op = "pgsql.GetOrder"
	var res models.GetOrdersResponse

	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		log.Printf("error from operation: %s: %s", op, err.Error())
		return res, fmt.Errorf("could not begin transaction: %v", err)
	}

	// getting all order_id's of this user (prepare)
	intOrders := make([]int, 0, 10)
	Orders := make([]models.Order, 0, 10)
	getOrderIDsByUserIDQuery := "SELECT * FROM orders WHERE user_id = $1"
	rows, err := o.db.pool.Query(ctx, getOrderIDsByUserIDQuery, userID)
	if err != nil {
		log.Printf("error from operation: %s: %s", op, err.Error())
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var or models.Order
		err := rows.Scan(&or.ID, &or.Time, &or.OrderStatus, &or.UserID, &or.PaymentMethod, &or.Cost)
		if err != nil {
			log.Printf("error from operation: %s: %s", op, err.Error())
			return res, fmt.Errorf("could not begin transaction: %v", err)
		}
		intOrders = append(intOrders, or.ID)
		Orders = append(Orders, or)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error from operation: %s: %s", op, err.Error())
		return res, err
	}

	// main loop
	for i, order_id := range intOrders {
		// cakes list creation
		intCakes := make([]int, 0, 10)
		getCakeIDQuery := "SELECT cake_id FROM order_cakes WHERE order_id = $1"
		rows, err := o.db.pool.Query(ctx, getCakeIDQuery, order_id)
		if err != nil {
			log.Printf("error from operation: %s: %s", op, err.Error())
			return res, err
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				log.Printf("error from operation: %s: %s", op, err.Error())
				return res, err
			}
			intCakes = append(intCakes, id)
		}
		if err := rows.Err(); err != nil {
			log.Printf("error from operation: %s: %s", op, err.Error())
			return res, err
		}

		// next station is cakes table!
		cakes := make([]models.Cake, 0, 10)
		for _, cakeID := range intCakes {
			getCakesQuery := "SELECT * FROM cakes WHERE id = $1"
			rows, err := o.db.pool.Query(ctx, getCakesQuery, cakeID)
			if err != nil {
				log.Printf("error from operation: %s: %s", op, err.Error())
				return res, err
			}
			defer rows.Close()
			for rows.Next() {
				var cake models.Cake
				err := rows.Scan(&cake.ID, &cake.Description, &cake.Price, &cake.Weight)
				if err != nil {
					log.Printf("error from operation: %s: %s", op, err.Error())
					return res, err
				}
				cakes = append(cakes, cake)
			}
			if err := rows.Err(); err != nil {
				log.Printf("error from operation: %s: %s", op, err.Error())
				return res, err
			}
		}

		ord := models.InternOrder{
			Cakes: cakes,
			Ord:   Orders[i],
		}
		res.Orders = append(res.Orders, ord)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("error from operation: %s: %s", op, err.Error())
		return res, err
	}

	return res, nil
}

func (o *UserOrderManagerPostgres) DeleteOrder(userID, orderID int) error {
	const op = "pgsql.DeleteOrder"

	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("op: %s: could not begin transaction: %s", op, err.Error())
	}

	// validation:
	query, args, err := sq.Select("user_id").
		From(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": orderID}).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't exec: %s: %s", op, err.Error())
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't exec: %s: %s", op, err.Error())
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("can't find user's order: %s: %s", op, "no rows affected")
	}

	query, args, err = sq.Delete(OrdersCakesTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s: %s", op, err.Error())
	}

	res, err = o.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s: %s", op, err.Error())
	}

	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("op: %s: %s", op, "could not delete order")
	}

	query, args, err = sq.Delete(DeliveryTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s: %s", op, err.Error())
	}

	res, err = o.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s: %s", op, err.Error())
	}
	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("op: %s: %s", op, "could not delete order")
	}

	query, args, err = sq.Delete(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s: %s", op, err.Error())
	}
	res, err = o.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s: %s", op, err.Error())
	}
	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("op: %s: %s", op, "could not delete order")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("op: %s: could not commit transaction: %s", op, err.Error())
	}

	return nil
}
