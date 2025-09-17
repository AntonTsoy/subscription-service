package main

import (
	"fmt"

	"github.com/AntonTsoy/subscription-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Конфиг: %+v\n", cfg)
}
