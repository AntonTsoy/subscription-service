package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
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
	r.Use(middleware.Timeout(10 * time.Second))

	r.Post("/subscriptions", subsHandler.CreateSubscription)
	r.Get("/subscriptions/{id}", subsHandler.GetSubscription)
	r.Get("/subscriptions", subsHandler.GetAllSubscriptions)
	r.Put("/subscriptions/{id}", subsHandler.UpdateSubscription)
	r.Delete("/subscriptions/{id}", subsHandler.DeleteSubscription)
	r.Get("/subscriptions/{start_date}/{end_date}/total-cost", subsHandler.TotalServiceSubscriptionsCost)
	r.Get("/subscriptions/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	fmt.Println("Listen on http://localhost:8080/subscriptions/hello")
	http.ListenAndServe(":8080", r)
}
