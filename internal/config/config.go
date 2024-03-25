package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	StorageData
}

type StorageData struct {
	User     string
	Password string
	DBName   string
}

func Load() (*Config, error) {
	var cfg Config

	cfg.User = os.Getenv("user")
	cfg.Password = os.Getenv("password")
	cfg.DBName = os.Getenv("dbName")

	if cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("environment variable is not set")
	}

	log.Println("OK: Config successfully loaded.")

	return &cfg, nil
}
