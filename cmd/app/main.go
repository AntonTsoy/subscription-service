package main

import (
	"log"

	"github.com/AntonTsoy/subscription-service/internal/config"
	"github.com/AntonTsoy/subscription-service/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка конфига: %v", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("ошибка базы данных: %v", err)
	}
	defer db.Close()

	if err = db.HealthCheck(); err != nil {
		log.Fatalf("не удалось открыть соединение c базой данных: %v", err)
	}
}
