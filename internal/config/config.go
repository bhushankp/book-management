package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Auth struct {
		JWTSecret string
	}
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// env overrides: APP_SERVER_PORT etc.
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	_ = viper.BindEnv("server.port", "APP_SERVER_PORT")
	_ = viper.BindEnv("database.host", "APP_DB_HOST")
	_ = viper.BindEnv("database.port", "APP_DB_PORT")
	_ = viper.BindEnv("database.user", "APP_DB_USER")
	_ = viper.BindEnv("database.password", "APP_DB_PASSWORD")
	_ = viper.BindEnv("database.name", "APP_DB_NAME")
	_ = viper.BindEnv("database.sslmode", "APP_DB_SSLMODE")
	_ = viper.BindEnv("auth.jwtsecret", "APP_JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}
	return &cfg
}
