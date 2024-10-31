package storage

import (
	"cakes-database-app/pkg/models"
	"context"
)

type CakeManagerPostgres struct {
	db *DB
}

func NewCakeManagerPostgres(db *DB) *CakeManagerPostgres {
	return &CakeManagerPostgres{db: db}
}

func (c *CakeManagerPostgres) GetCakes() ([]models.Cake, error) {
	cakes := make([]models.Cake, 0, 10)
	query := "SELECT id, description, price, weight FROM cakes"
	rows, err := c.db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cake models.Cake
		err := rows.Scan(&cake.ID, &cake.Description, &cake.Price, &cake.Weight)
		if err != nil {
			return nil, err
		}
		cakes = append(cakes, cake)
	} 
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cakes, nil
}