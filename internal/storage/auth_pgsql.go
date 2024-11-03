package storage

import (
	"cakes-database-app/internal/models"
	"context"
)

type AuthPostgres struct {
	db *DB
}

func NewAuthPostgres(db *DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user models.User) (int, error) {
	var userID int
	err := a.db.pool.QueryRow(context.Background(),
		"INSERT INTO users (fullname, username, email, password_hash, phone_number) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.PhoneNumber).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *AuthPostgres) GetUser(username, password_hash string) (int, error) {
	var userID int
	err := a.db.pool.QueryRow(context.Background(),
		"SELECT id FROM users WHERE username = $1 AND password_hash = $2",
		username, password_hash).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
