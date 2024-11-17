package service

import (
	"context"
	"errors"

	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"

	"github.com/golang-jwt/jwt/v5"
)

const adminSigningKey = "a1kvoai45is9cklaa;sgdhsudvk"

var _ AdminAuthorization = (*AdminAuthorizationService)(nil)
var _ AdminService = (*AdminServ)(nil)

type AdminAuthorizationService struct {
	stg storage.AdminAuthorization
}

type AdminServ struct {
	stg storage.Admin
}

func NewAdminAuthService(stg storage.AdminAuthorization) *AdminAuthorizationService {
	return &AdminAuthorizationService{stg: stg}
}

type adminTokenClaims struct {
	jwt.RegisteredClaims
	AdminID int `json:"admin_id"`
}

func (s *AdminAuthorizationService) GenerateAdminToken(username, password string) (string, error) {
	// get user from db
	adminID, err := s.stg.GetAdmin(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &adminTokenClaims{
		AdminID: adminID,
	})

	return token.SignedString([]byte(adminSigningKey))
}

func (s *AdminAuthorizationService) ParseAdminToken(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &adminTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(adminSigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*adminTokenClaims)
	if !ok || claims == nil {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.AdminID, nil
}

func NewAdminService(stg storage.Admin) *AdminServ {
	return &AdminServ{stg: stg}
}

func (a *AdminServ) GetUsers() ([]models.User, error) {
	return a.stg.GetUsers()
}

func (a *AdminServ) Backup() error {
	return a.stg.Backup()
}

func (a *AdminServ) Restore() error {
	return a.stg.Restore()
}

func (a *AdminServ) AddCake(ctx context.Context, cake models.Cake) (int, error) {
	return a.stg.AddCake(ctx, cake)
}

func (a *AdminServ) RemoveCake(ctx context.Context, cake int) error {
	return a.stg.RemoveCake(ctx, cake)
}
