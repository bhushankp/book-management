package main

import (
	"book-management/internal/config"
	"book-management/internal/db"
	"book-management/internal/models"
	"book-management/internal/pkg/logger"
	"book-management/internal/pkg/validator"
	"fmt"

	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadConfig()
	logger.InitLogger()
	db.Connect(cfg)
	validator.Init()
	logger.Log.Info("config loaded",
		zap.Int("port", cfg.Server.Port),
	)

	book := models.Book{Title: "", Author: "Unknown"}
	if err := validator.ValidateStruct(book); err != nil {
		fmt.Println("Validation failed:", err)
	}

}
