package pgsql

import (
	"cakes-database-app/pkg/models"
	"context"

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


func (db *DB) GetCakes(ctx context.Context) ([]models.Cake, error) {
	// TODO: implement this
	return nil, nil
}
