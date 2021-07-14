package lib

import (
	"go.uber.org/zap"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once

	HttpServerLog = zap.String("module", "http-server")
	ChromeLog     = zap.String("module", "chrome")
)

func Logger() *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()

		if err != nil {
			panic(err)
		}
	})

	return logger
}
