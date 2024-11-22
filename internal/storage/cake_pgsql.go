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
	const op = "pgsql.GetCakes"
	cakes := make([]models.Cake, 0, 50)

	builderSelect := sq.Select(idColumn, descriptionColumn, priceColumn, weightColumn, fullDescriptionColumn).
		From(CakesTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{activeColumn: true})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return cakes, fmt.Errorf("op: %s, failed to build query: %w", op, err)
	}

	rows, err := c.db.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cake models.Cake
		err := rows.Scan(&cake.ID, &cake.Description, &cake.Price, &cake.Weight, &cake.FullDescription)
		if err != nil {
			return nil, fmt.Errorf("op: %s, failed to scan row: %w", op, err)
		}
		cakes = append(cakes, cake)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("op: %s, failed to scan rows: %w", op, err)
	}

	return cakes, nil
}
