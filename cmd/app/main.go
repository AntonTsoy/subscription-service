package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AntonTsoy/subscription-service/internal/config"
	"github.com/AntonTsoy/subscription-service/internal/database"
	"github.com/AntonTsoy/subscription-service/internal/repository"
	"github.com/AntonTsoy/subscription-service/internal/service"
	"github.com/AntonTsoy/subscription-service/internal/transport/handler"
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

	subsRepo := repository.NewSubsRepo(db.DB())

	subsService := service.NewSubsService(subsRepo)

	subsHandler := handler.NewSubsHandler(subsService)

	r := chi.NewRouter()
	r.Post("/subscriptions", subsHandler.CreateSubscription)
	r.Get("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	fmt.Println("Listen on http://localhost:8080/subscriptions")
	http.ListenAndServe(":8080", r)
}
