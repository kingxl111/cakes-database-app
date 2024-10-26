package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// entities
const (
	userTable = 			"users"
	deliveryTable = 		"delivery"
	orderTable = 			"orders"
	ordersCakesTable = 		"orders_cakes"
	cakesTable = 			"cakes"
	deliveryPointTable = 	"delivery_point"
	adminTable = 			"administrators"
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

