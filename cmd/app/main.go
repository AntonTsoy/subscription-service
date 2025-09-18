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

	//subsRepo := repository.NewSubsRepo(db.DB())

	/*
		userId, _ := uuid.Parse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
		layout := "01-2006"
		startDate, _ := time.Parse(layout, "07-2025")
		testModel := models.Subscription{ServiceName: "Spotify", Price: 20, UserID: userId, StartDate: startDate}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		subsRepo.Create(ctx, &testModel)

		fmt.Printf("%v", testModel)
	*/
}
