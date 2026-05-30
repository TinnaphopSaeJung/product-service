package main

import (
	"context"
	"log"
	"time"

	"product-service/internal/app"
	"product-service/internal/config"
	"product-service/internal/database"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer dbPool.Close()

	router := app.NewRouter(dbPool)

	log.Printf("Server is running on port %s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
