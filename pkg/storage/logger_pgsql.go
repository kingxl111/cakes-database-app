package storage

import (
	"context"
	"time"
)

type LoggerPostgres struct {
	db *DB
}

func NewLoggerPostgres(db *DB) *LoggerPostgres {
	return &LoggerPostgres{db: db}
}

func (l *LoggerPostgres) WriteLog(lvl string, msg string) error {
	_, err := l.db.pool.Exec(context.Background(), "INSERT INTO logs (level, message, created_at) VALUES ($1, $2, $3)", lvl, msg, time.Now())
    return err
}	