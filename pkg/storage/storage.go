package storage

import (
	"cakes-database-app/pkg/models"
	"context"
)

type Interface interface {
	GetCakes(context.Context) ([]models.Cake, error)
}

