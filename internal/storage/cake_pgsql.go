package storage

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/kingxl111/cakes-database-app/internal/models"
)

type UserCakeManagerPostgres struct {
	db *DB
}

func NewUserCakeManagerPostgres(db *DB) *UserCakeManagerPostgres {
	return &UserCakeManagerPostgres{db: db}
}

func (c *UserCakeManagerPostgres) GetCakes() ([]models.Cake, error) {
	cakes := make([]models.Cake, 0, 10)

	builderSelect := sq.Select("id", "description", "price", "weight", "full_description").
		From(CakesTable).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return cakes, fmt.Errorf("failed to build query: %v", err.Error())
	}

	rows, err := c.db.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cake models.Cake
		err := rows.Scan(&cake.ID, &cake.Description, &cake.Price, &cake.Weight, &cake.FullDescription)
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
