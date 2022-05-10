package logging

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() (*zap.SugaredLogger, func(), error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		errSync := logger.Sync()
		if errSync != nil {
			log.Printf("Cannot flush any buffer log entries")
			return
		}
	}

	return logger.Sugar(), cleanup, nil
}
