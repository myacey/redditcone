package logging

import (
	"os"

	"go.uber.org/zap"
)

func ConfigureLogger() *zap.SugaredLogger {
	logger := zap.Must(zap.NewDevelopment())
	if os.Getenv("LOGGER_TYPE") == "production" {
		logger = zap.Must(zap.NewProduction())
	}

	sugar := logger.Sugar()
	sugar.Info("logger initialized")

	return sugar
}
