package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
	JWTDuration time.Duration
	BcryptCost  int
}

func Load() *Config {
	_ = godotenv.Load()

	duration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		log.Fatal("invalid JWT_DURATION")
	}

	return &Config{
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTDuration: duration,
		BcryptCost:  12,
	}
}
