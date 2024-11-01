package service

import (
	"cakes-database-app/pkg/storage"
)

type LoggerService struct {
	stg storage.Logger
}

func NewLoggerService(stg storage.Logger) *LoggerService {
	return &LoggerService{stg: stg}
} 

func (l *LoggerService)WriteLog(level string, msg string) error {
	return l.stg.WriteLog(level, msg)
}