package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/kingxl111/cakes-database-app/internal/storage/s3"
	"github.com/sirupsen/logrus"

	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
)

type CakeService struct {
	stg  storage.UserCakeManager
	serv s3.ClientS3
}

func NewCakeService(stg storage.UserCakeManager, serv s3.ClientS3) *CakeService {
	return &CakeService{
		stg:  stg,
		serv: serv,
	}
}

func (c *CakeService) GetCakes() ([]models.Cake, error) {
	cakes, err := c.stg.GetCakes()
	if err != nil {
		return cakes, err
	}
	for i, _ := range cakes {
		imageUrl := c.serv.GetFileURL(cakes[i].Description)
		cakes[i].ImageURL = imageUrl
		logrus.Printf("cake: %s, image_url: %s\n", cakes[i].Description, cakes[i].ImageURL)
	}
	return cakes, err
}

func (c *CakeService) GetCake(id int) (models.Cake, error) {
	ar, err := c.GetCakes()
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

func (c *CakeService) UploadCakePhoto(ctx context.Context, file multipart.File, fileName string) (string, error) {
	return c.serv.UploadFile(ctx, file, fileName)
}
