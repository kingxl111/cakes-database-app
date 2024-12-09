package service

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/kingxl111/cakes-database-app/internal/models"
	"github.com/kingxl111/cakes-database-app/internal/storage"

	"github.com/golang-jwt/jwt/v5"
)

const (
	salt       = "jkaawken11elzc;d12k2Wfpcallsdhac" // for hash
	signingKey = "kwoduehcziweligfj29kxz.8ck"       // for jwt signing
)

type tokenClaims struct {
	jwt.RegisteredClaims     // Claims standards
	UserId               int `json:"userId"`
}

type AuthService struct {
	stg storage.Authorization
}

func NewAuthService(stg storage.Authorization) *AuthService {
	return &AuthService{stg: stg}
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.PasswordHash = generatePasswordHash(user.PasswordHash)
	//log.Printf("hash: %s", user.PasswordHash)
	// other user's fields without changes
	return s.stg.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	// get user from db
	userID, err := s.stg.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		UserId: userID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || claims == nil {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
