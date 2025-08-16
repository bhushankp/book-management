package main

import (
	"book-management/internal/config"
	"book-management/internal/pkg/logger"

	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadConfig()
	logger.InitLogger()

	logger.Log.Info("config loaded",
		zap.Int("port", cfg.Server.Port),
	)

}
