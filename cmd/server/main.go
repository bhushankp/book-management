package main

import (
	"book-management/internal/config"
	"book-management/internal/db"
	"book-management/internal/pkg/logger"

	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadConfig()
	logger.InitLogger()
	db.Connect(cfg)

	logger.Log.Info("config loaded",
		zap.Int("port", cfg.Server.Port),
	)

}
