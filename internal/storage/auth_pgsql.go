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

func (a *AuthPostgres) CreateUser(user models.User) (int, error) {
	const op = "pgsql.CreateUser"
	var userID int

	builderInsert := sq.Insert(UserTable).
		PlaceholderFormat(sq.Dollar).
		Columns(fullnameColumn, usernameColumn, emailColumn, passwordColumn, phoneColumn).
		Values(user.FullName, user.Username, user.Email, user.PasswordHash, user.PhoneNumber).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return userID, fmt.Errorf("op: %s, failed to build query: %w", op, err)
	}

	err = a.db.pool.QueryRow(context.Background(), query, args...).Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("op: %s, failed to execute query: %w", op, err)
	}
	return userID, nil
}

func (a *AuthPostgres) GetUser(username, passwordHash string) (int, error) {
	const op = "pgsql.GetUser"
	var userID int

	builder := sq.Select(idColumn).
		From(UserTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{usernameColumn: username}).
		Where(sq.Eq{passwordColumn: passwordHash})

	query, args, err := builder.ToSql()
	if err != nil {
		return userID, fmt.Errorf("op: %s, failed to build query: %w", op, err)
	}

	err = a.db.pool.QueryRow(context.Background(), query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
