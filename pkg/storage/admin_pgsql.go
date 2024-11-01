package storage

import (
	"cakes-database-app/pkg/models"
	"context"
)

type AdminAuthPostgres struct {
	db *DB
}

func NewAdminAuthPostgres(db *DB) *AdminAuthPostgres {
	return &AdminAuthPostgres{db: db}
}

func (a *AdminAuthPostgres) GetAdmin(username, password_hash string) (int, error) {
	var id int
	ctx := context.Background()
	query := "SELECT id FROM admins WHERE username = $1 AND password_hash = $2"
	err := a.db.pool.QueryRow(ctx, query, username, password_hash).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

type AdminPostgres struct {
	db *DB
}

func NewAdminPostgres(db *DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (a *AdminPostgres) GetUsers() ([]models.User, error) {
	users := make([]models.User, 0, 10)
	query := "SELECT id, fullname, username, email, phone_number FROM users"
	rows, err := a.db.pool.Query(context.Background(), query)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var usr models.User
		err = rows.Scan(&usr.ID, &usr.FullName, &usr.Username, &usr.Email, &usr.PhoneNumber)
		if err != nil {
			return users, err
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}
