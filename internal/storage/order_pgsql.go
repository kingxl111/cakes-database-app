package storage

import (
	"context"
	"fmt"
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
	const op = "pgsql.CreateOrder"
	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("op: %s, could not begin transaction: %w", op, err)
	}

	// calculate common order's cakes cost by cakes ids:
	cost := 0
	for _, cake := range cakes {
		var price int
		id := cake.ID

		builderSelect := sq.Select(priceColumn).
			From(CakesTable).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{activeColumn: true}).
			Where(sq.Eq{idColumn: id})

		query, args, err := builderSelect.ToSql()
		if err != nil {
			return orderID, fmt.Errorf("op: %s, could not build query: %w", op, err)
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
		Columns(timestampColumn, orderStatusColumn, userIdColumn, paymentMethodColumn, costColumn).
		Values(order.Time, order.OrderStatus, order.UserID, order.PaymentMethod, order.Cost).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		return orderID, fmt.Errorf("op: %s, could not build query: %w", op, err)
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&orderID)
	if err != nil {
		return orderID, fmt.Errorf("op: %s, could not insert order: %w", op, err)
	}

	var deliveryID int // the same with orderID
	builderInsert = sq.Insert(DeliveryTable).
		PlaceholderFormat(sq.Dollar).
		Columns(pointIdColumn, costColumn, deliveryStatusColumn, weightColumn).
		Values(delivery.PointID, delivery.Cost, delivery.Status, delivery.Weight).
		Suffix("RETURNING id")

	query, args, err = builderInsert.ToSql()
	if err != nil {
		return orderID, fmt.Errorf("op: %s, could not build query: %w", op, err)
	}
	err = tx.QueryRow(ctx, query, args...).Scan(&deliveryID)
	if err != nil {
		return orderID, fmt.Errorf("op: %s, could not insert delivery: %w", op, err)
	}

	for _, cake := range cakes {
		orderCake := models.OrderCake{
			OrderID: orderID,
			CakeID:  cake.ID,
		}
		builder := sq.Insert(OrdersCakesTable).
			PlaceholderFormat(sq.Dollar).
			Columns(orderIdColumn, cakeIdColumn).
			Values(orderCake.OrderID, orderCake.CakeID)
		query, args, err := builder.ToSql()
		if err != nil {
			return orderID, fmt.Errorf("op: %s, could not build query: %w", op, err)
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return orderID, fmt.Errorf("op: %s, could not insert order: %w", op, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return orderID, fmt.Errorf("op: %s, could not commit transaction: %w", op, err)
	}

	return orderID, nil
}

func (o *UserOrderManagerPostgres) GetOrders(userID int) (models.GetOrdersResponse, error) {
	const op = "pgsql.GetOrder"
	var res models.GetOrdersResponse

	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return res, fmt.Errorf("op: %s, could not begin transaction: %w", op, err)
	}

	// getting all order_id's of this user (prepare)
	intOrders := make([]int, 0, 10)
	Orders := make([]models.Order, 0, 10)

	builder := sq.Select("*").
		From(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIdColumn: userID}).
		Where(sq.NotEq{orderStatusColumn: "canceled"})
	query, args, err := builder.ToSql()
	if err != nil {
		return res, fmt.Errorf("op: %s, could not build query: %w", op, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return res, fmt.Errorf("op: %s, could not query rows: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var or models.Order
		err := rows.Scan(&or.ID, &or.Time, &or.OrderStatus, &or.UserID, &or.PaymentMethod, &or.Cost)
		if err != nil {
			return res, fmt.Errorf("op: %s, could not scan row: %w", op, err)
		}
		intOrders = append(intOrders, or.ID)
		Orders = append(Orders, or)
	}
	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("op: %s, could not iterate rows: %w", op, err)
	}

	// main loop
	for i, orderId := range intOrders {
		// cakes list creation
		intCakes := make([]int, 0, 10)
		builder := sq.Select(cakeIdColumn).
			From(OrdersCakesTable).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{orderIdColumn: orderId})
		query, args, err := builder.ToSql()
		if err != nil {
			return res, fmt.Errorf("op: %s, could not build select query: %w", op, err)
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return res, fmt.Errorf("op: %s, could not query row: %w", op, err)
		}

		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				return res, fmt.Errorf("op: %s, could not scan row: %w", op, err)
			}
			intCakes = append(intCakes, id)
		}
		if err := rows.Err(); err != nil {
			return res, fmt.Errorf("op: %s, could not iterate rows: %w", op, err)
		}

		// next station is cakes table!
		cakes := make([]models.Cake, 0, 10)
		for _, cakeID := range intCakes {
			builder := sq.Select(idColumn, descriptionColumn, priceColumn, weightColumn, fullDescriptionColumn, activeColumn).
				From(CakesTable).
				PlaceholderFormat(sq.Dollar).
				Where(sq.Eq{activeColumn: true}).
				Where(sq.Eq{idColumn: cakeID})
			query, args, err := builder.ToSql()
			if err != nil {
				return res, fmt.Errorf("op: %s, could not build select query: %w", op, err)
			}

			rows, err := tx.Query(ctx, query, args...)
			if err != nil {
				return res, fmt.Errorf("op: %s, could not query row: %w", op, err)
			}

			for rows.Next() {
				var cake models.Cake
				err := rows.Scan(&cake.ID, &cake.Description, &cake.Price, &cake.Weight, &cake.FullDescription, &cake.Active)
				if err != nil {
					return res, fmt.Errorf("op: %s, could not scan row: %w", op, err)
				}
				cakes = append(cakes, cake)
			}
			if err := rows.Err(); err != nil {
				return res, fmt.Errorf("op: %s, could not iterate rows: %w", op, err)
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
		return res, fmt.Errorf("op: %s, could not commit transaction: %w", op, err)
	}

	return res, nil
}

func (o *UserOrderManagerPostgres) DeleteOrder(userID, orderID int) error {
	const op = "pgsql.DeleteOrder"

	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("op: %s, could not begin transaction: %w", op, err)
	}

	// validation:
	query, args, err := sq.Select(userIdColumn).
		From(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: orderID}).
		Where(sq.Eq{userIdColumn: userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s, could not build query: %w", op, err)
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s, could not query row: %w", op, err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("op: %s, no rows affected: %w", op, err)
	}

	query, args, err = sq.Update(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Set(orderStatusColumn, "canceled").
		Where(sq.Eq{idColumn: orderID}).
		Where(sq.Eq{userIdColumn: userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s, could not build update query: %w", op, err)
	}
	res, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s, could not query row: %w", op, err)
	}
	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("op: %s, no rows affected: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("op: %s, could not commit transaction: %w", op, err)
	}

	return nil
}

func (o *UserOrderManagerPostgres) UpdateOrder(userID int, orderID int, paymentMethod string) error {
	const op = "pgsql.UpdateOrder"
	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("op: %s, could not begin transaction: %w", op, err)
	}

	query, args, err := sq.Update(OrderTable).
		PlaceholderFormat(sq.Dollar).
		Set(paymentMethodColumn, paymentMethod).
		Where(sq.Eq{idColumn: orderID}).
		Where(sq.Eq{userIdColumn: userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s, could not build update query: %w", op, err)
	}
	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s, could not query row: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("op: %s, no rows affected: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("op: %s, could not commit transaction: %w", op, err)
	}
	return nil
}

func (o *UserOrderManagerPostgres) GetDeliveryPoints() ([]models.DeliveryPoint, error) {
	const op = "pgsql.GetDeliveryPoints"

	points := make([]models.DeliveryPoint, 0, 10)
	ctx := context.Background()
	tx, err := o.db.pool.Begin(ctx)
	if err != nil {
		return points, fmt.Errorf("op: %s, could not begin transaction: %w", op, err)
	}

	qeury, args, err := sq.Select(idColumn, ratingColumn, addressColumn, workingHoursColumn, contactPhoneColumn).
		From(DeliveryPointTable).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return points, fmt.Errorf("op: %s, could not build query: %w", op, err)
	}
	res, err := tx.Query(ctx, qeury, args...)
	if err != nil {
		return points, fmt.Errorf("op: %s, could not query row: %w", op, err)
	}

	for res.Next() {
		var point models.DeliveryPoint
		err := res.Scan(&point.ID, &point.Rating, &point.Address, &point.WorkingHours, &point.ContactPhone)
		if err != nil {
			return points, fmt.Errorf("op: %s, could not scan row: %w", op, err)
		}
		points = append(points, point)
	}
	if err := res.Err(); err != nil {
		return points, fmt.Errorf("op: %s, could not iterate rows: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return points, fmt.Errorf("op: %s, could not commit transaction: %w", op, err)
	}

	return points, nil
}
