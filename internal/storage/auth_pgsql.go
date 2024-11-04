package storage

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kingxl111/cakes-database-app/internal/models"
)

type AuthPostgres struct {
	db *DB
}

func NewAuthPostgres(db *DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

/*
type User struct {
	ID           int    `json:"id"`
	FullName     string `json:"fullname"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	PhoneNumber  string `json:"phone_number"`
}
*/

func (a *AuthPostgres) CreateUser(user models.User) (int, error) {
	var userID int

	builderInsert := sq.Insert(UserTable).
		PlaceholderFormat(sq.Dollar).
		Columns("fullname", "username", "email", "password_hash", "phone_number").
		Values(user.FullName, user.Username, user.Email, user.PasswordHash, user.PhoneNumber).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return userID, fmt.Errorf("failed to build query: %v", err.Error())
	}

	err = a.db.pool.QueryRow(context.Background(), query, args...).Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("failed to insert user: %v", err.Error())
	}
	return userID, nil
}

func (a *AuthPostgres) GetUser(username, passwordHash string) (int, error) {
	var userID int

	builder := sq.Select("id").
		From(UserTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"username": username}).
		Where(sq.Eq{"password_hash": passwordHash})

	query, args, err := builder.ToSql()
	if err != nil {
		return userID, fmt.Errorf("failed to build query: %v", err.Error())
	}

	err = a.db.pool.QueryRow(context.Background(), query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
