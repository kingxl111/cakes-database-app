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
	DeliveryPointTable = "delivery_point"
	AdminTable         = "administrators"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(username, password, address, dbname, sslmode string) (*DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, address, dbname, sslmode)
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
