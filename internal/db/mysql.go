package db

import (
	"fmt"

	"book-management/internal/config"
	"book-management/internal/models"
	"book-management/internal/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("db connect failed", zap.Error(err))
	}
	if err := DB.AutoMigrate(&models.Book{}); err != nil {
		logger.Log.Fatal("db migrate failed", zap.Error(err))
	}
	logger.Log.Info("db connected")
}
