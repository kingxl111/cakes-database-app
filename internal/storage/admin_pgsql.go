package storage

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/kingxl111/cakes-database-app/internal/models"
)

type AdminAuthPostgres struct {
	db *DB
}

func NewAdminAuthPostgres(db *DB) *AdminAuthPostgres {
	return &AdminAuthPostgres{db: db}
}

func (a *AdminAuthPostgres) GetAdmin(username, passwordHash string) (int, error) {
	var id int
	//log.Println("password_hash: " + passwordHash)
	builderSelect := sq.Select("id").
		From(AdminTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"username": username}).
		Where(sq.Eq{"password_hash": passwordHash})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return id, fmt.Errorf("error building query: %v", err.Error())
	}
	log.Println(args...)
	err = a.db.pool.QueryRow(context.Background(), query, args...).Scan(&id)
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

	builderSelect := sq.Select("id", "fullname", "username", "email", "phone_number").
		From(UserTable).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return users, fmt.Errorf("error building query: %v", err.Error())
	}

	rows, err := a.db.pool.Query(context.Background(), query, args...)
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

func (a *AdminPostgres) Backup() error {

	return nil
}

//func (a *AdminPostgres) DeleteUser(userID int) error {
//	builderDelete := sq.Delete(UserTable).
//		PlaceholderFormat(sq.Dollar).
//		Where(sq.Eq{"id": userID})
//
//	query, args, err := builderDelete.ToSql()
//	if err != nil {
//		return fmt.Errorf("error building query: %v", err.Error())
//	}
//
//	res, err := a.db.pool.Exec(context.Background(), query, args...)
//	if err != nil {
//		return fmt.Errorf("error deleting user: %v", err.Error())
//	}
//
//	return nil
//}
