package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"

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
	redis    *redis.Client
	stg      storage.Authorization
	tokenTTL time.Duration
}

func NewAuthService(stg storage.Authorization, rdb *redis.Client, ttl time.Duration) *AuthService {
	return &AuthService{stg: stg, redis: rdb, tokenTTL: ttl}
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
	userID, err := s.stg.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	s.tokenTTL = time.Minute * 15
	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
		},
		UserId: userID,
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("Token ExpiresAt:", claims.ExpiresAt)
	tokenStr, err := tokenObj.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	log.Printf("DEBUG: setting redis key %q → %d (TTL %s)\n", "token:"+tokenStr, userID, s.tokenTTL)
	//key := fmt.Sprintf("token:%s", tokenStr)
	key := "token:" + tokenStr
	ctx := context.Background()
	if err := s.redis.Set(ctx, key, userID, s.tokenTTL).Err(); err != nil {
		return "", err
	}
	n, _ := s.redis.Exists(ctx, key).Result()
	log.Printf("DEBUG: Exists after SET: %d (ожидается 1)\n", n)

	return tokenStr, nil
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
		return 0, errors.New("invalid token claims")
	}

	key := "token:" + accessToken
	log.Printf("DEBUG: looking up redis key %q\n", key)
	//key := fmt.Sprintf("token:%s", accessToken)
	val, err := s.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, errors.New("token expired or not found")
	}
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}
