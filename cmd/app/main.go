package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/AntonTsoy/subscription-service/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/AntonTsoy/subscription-service/internal/config"
	"github.com/AntonTsoy/subscription-service/internal/database"
	"github.com/AntonTsoy/subscription-service/internal/repository"
	"github.com/AntonTsoy/subscription-service/internal/service"
	"github.com/AntonTsoy/subscription-service/internal/transport/handler"
	"github.com/AntonTsoy/subscription-service/internal/transport/logger"
)

// @title           Subscription Service API
// @version         1.0
// @description     REST API для управления подписками
// @BasePath        /
// @host            localhost:8080
// @schemes         http
func main() {
	cfg := config.Load()

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
	r.Use(logger.Logger)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Post("/subscriptions", subsHandler.CreateSubscription)
	r.Get("/subscriptions/{id}", subsHandler.GetSubscription)
	r.Get("/subscriptions", subsHandler.GetAllSubscriptions)
	r.Put("/subscriptions/{id}", subsHandler.UpdateSubscription)
	r.Delete("/subscriptions/{id}", subsHandler.DeleteSubscription)
	r.Get("/subscriptions/{start}/{end}/total-cost", subsHandler.TotalServiceSubscriptionsCost)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Server started at http://localhost:8080/swagger/index.html")
	http.ListenAndServe(":8080", r)
}
