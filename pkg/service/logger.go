package service

import (
	"cakes-database-app/pkg/storage"
	"context"
)

type LoggerService struct {
	stg storage.Logger
}

func NewLoggerService(stg *storage.Storage) *LoggerService {
	return &LoggerService{stg: stg}
} 

func (l *LoggerService)WriteLog(ctx *context.Context, level string, msg string) error {
	return l.stg.WriteLog(ctx, level, msg)
}