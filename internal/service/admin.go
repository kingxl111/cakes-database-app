package service

import (
	"cakes-database-app/internal/models"
	"cakes-database-app/internal/storage"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const adminSigningKey = "a1kvoai45is9cklaa;sgdhsudvk"

type AdminAuthorizationService struct {
	stg storage.AdminAuthorization
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

type AdminServ struct {
	stg storage.Admin
}

func NewAdminService(stg storage.Admin) *AdminServ {
	return &AdminServ{stg: stg}
}

func (a *AdminServ) GetUsers() ([]models.User, error) {
	return a.stg.GetUsers()
}
