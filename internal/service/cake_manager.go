package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mime/multipart"
	"time"

	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"
	"github.com/kingxl111/cakes-database-app/internal/storage/s3"
)

type CakeService struct {
	stg      storage.UserCakeManager
	serv     s3.ClientS3
	redis    *redis.Client
	cacheTTL time.Duration
}

func NewCakeService(
	stg storage.UserCakeManager,
	serv s3.ClientS3,
	rdb *redis.Client,
	cacheTTL time.Duration,
) *CakeService {
	return &CakeService{
		stg:      stg,
		serv:     serv,
		redis:    rdb,
		cacheTTL: cacheTTL,
	}
}
func (c *CakeService) GetCakes() ([]models.Cake, error) {
	ctx := context.Background()
	cacheKey := "cache:cakes"

	if data, err := c.redis.Get(ctx, cacheKey).Result(); err == nil {
		var cakes []models.Cake
		if err := json.Unmarshal([]byte(data), &cakes); err == nil {
			return cakes, nil
		}
	}

	cakes, err := c.stg.GetCakes()
	if err != nil {
		return nil, err
	}

	for i := range cakes {
		cakes[i].ImageURL = c.serv.GetFileURL(cakes[i].Description)
	}

	b, err := json.Marshal(cakes)
	if err == nil {
		_ = c.redis.Set(ctx, cacheKey, b, c.cacheTTL).Err()
	}

	return cakes, nil
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
