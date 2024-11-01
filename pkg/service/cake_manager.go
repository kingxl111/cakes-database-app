package service

import (
	"cakes-database-app/pkg/models"
	"cakes-database-app/pkg/storage"
)

type CakeService struct {
	stg storage.UserCakeManager
}

func NewCakeService(stg storage.UserCakeManager) *CakeService {
	return &CakeService{stg: stg}
} 

func (c *CakeService) GetCakes() ([]models.Cake, error) {
	return c.stg.GetCakes()
}