package pgsql

import (
	"cakes-database-app/pkg/storage"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}


func (db *DB) GetCakes(ctx context.Context) ([]storage.Cake, error) {
	// TODO: implement this
	return nil, nil
}

func (db *DB) AddCake(ctx context.Context, cake storage.Cake) error {
	// TODO: implement this
	return nil
}
