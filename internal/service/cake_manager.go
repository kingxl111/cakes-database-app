package service

import (
	"fmt"

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

func (c *CakeService) GetCake(id int) (models.Cake, error) {
	ar, err := c.stg.GetCakes()
	if err != nil {
		return models.Cake{}, err
	}
	for idx, cake := range ar {
		if cake.ID == id {
			return ar[idx], nil
		}
	}
	return models.Cake{}, fmt.Errorf("wrong cake index: %v", id)
}
