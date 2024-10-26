package service

import (
	"cakes-database-app/pkg/models"
	"cakes-database-app/pkg/storage"
	"crypto/sha1"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

const (
	salt = "jkaawken11elzc;d12k2Wfpcallsdhac"  // for hash
	signingKey = "skdf12njx83xu39ck2d91cwjel"  // for jwt signing
)

type tokenClaims struct {
	jwt.RegisteredClaims   		// Claims standards
	UserId int `json:"userId"`
}

type AuthService struct {
	stg storage.Authorization
}

func NewAuthService(stg storage.Authorization) *AuthService {
	return &AuthService{stg: stg}
}

// unexportable method
func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s* AuthService) CreateUser(user models.User) (int, error) {
	user.PasswordHash = s.generatePasswordHash(user.PasswordHash)
	log.Printf("hash: %s", user.PasswordHash)
	// other user's fields without changes
	return s.stg.CreateUser(user)
}

// func (s* AuthService) GenerateToken(username, password string) (string, error) {
// 	// get user from db
// 	user, err := s.stg.GetUser(username, s.generatePasswordHash(password))

// 	if err != nil {
// 		return "", err
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
// 		UserId: user.Id,
// 	})
	
// 	return token.SignedString([]byte(signingKey))
// }

// func (s *AuthService) ParseToken(accessToken string) (int, error) {

// 	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("Invalid signing method")
// 		}
// 		return []byte(signingKey), nil
// 	})
	

// 	if err != nil {
// 		return 0, err
// 	}

// 	claims, ok := token.Claims.(*tokenClaims)
// 	if !ok || claims == nil {
// 		return 0, errors.New("token claims are not of type *tokenClaims")
// 	}

// 	return claims.UserId, nil
// }

