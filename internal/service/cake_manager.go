package service

import (
	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
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
