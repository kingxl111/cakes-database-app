package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// entities
const (
	UserTable          = "users"
	DeliveryTable      = "deliveries"
	OrderTable         = "orders"
	OrdersCakesTable   = "order_cakes"
	CakesTable         = "cakes"
	DeliveryPointTable = "delivery_points"
	AdminTable         = "admins"

	idColumn       = "id"
	usernameColumn = "username"
	passwordColumn = "password_hash"
	emailColumn    = "email"
	phoneColumn    = "phone_number"
	fullnameColumn = "fullname"

	costColumn          = "cost"
	orderStatusColumn   = "order_status"
	paymentMethodColumn = "payment_method"
	userIdColumn        = "user_id"
	orderIdColumn       = "order_id"
	timestampColumn     = "time"

	addressColumn      = "address"
	ratingColumn       = "rating"
	workingHoursColumn = "working_hours"
	contactPhoneColumn = "contact_phone"

	pointIdColumn        = "point_id"
	deliveryStatusColumn = "status"
	weightColumn         = "weight"

	descriptionColumn     = "description"
	fullDescriptionColumn = "full_description"
	priceColumn           = "price"
	activeColumn          = "active"
	cakeIdColumn          = "cake_id"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(username, password, host, port, dbname, sslmode string) (*DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, host, port, dbname, sslmode)
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
