package service

import (
	"cakes-database-app/pkg/models"
	"cakes-database-app/pkg/storage"
)

type CakeService struct {
	stg storage.CakeManager
}

func NewCakeService(stg *storage.Storage) *CakeService {
	return &CakeService{stg: stg}
} 

func (c *CakeService) GetCakes() ([]models.Cake, error) {
	return c.stg.GetCakes()
}