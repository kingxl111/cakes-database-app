package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Lg *logrus.Logger
}

func NewLogger(logFilePath string) (*Logger, error) {
	log := logrus.New()

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(logrus.InfoLevel)

	return &Logger{Lg: log}, nil
}
