package pgsql

import (
	"cakes-database-app/pkg/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// entities
const (
	userTable = 			"users"
	deliveryTable = 		"delivery"
	orderTable = 			"orders"
	ordersCakesTable = 		"orders_cakes"
	cakesTable = 			"cakes"
	deliveryPointTable = 	"delivery_point"
	adminTable = 			"administrators"
)

type DB struct {
	pool *pgxpool.Pool
}


func NewDB(username, password, address, dbname, sslmode string) (*DB, error) {
    connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, address, dbname, sslmode)

    pool, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        return nil, err	
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, err
    }

    return &DB{pool: pool}, nil
}

func (db *DB) Close() {
    db.pool.Close()
}

func (db *DB) CreateUser(user models.User) (int, error) {
    var userID int
    err := db.pool.QueryRow(context.Background(),
        "INSERT INTO users (fullname, username, email, password_hash, phone_number) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        user.FullName,
		user.Username, 
		user.PasswordHash,
		user.Email,
		user.PhoneNumber).Scan(&userID)
    if err != nil {
        return 0, err
    }
    return userID, nil
}