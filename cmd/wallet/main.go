package main

import (
	"context"
	"log"
	"wallet/internal/app"
	"wallet/internal/config"
)

func main() {
	log.Println("Starting wallet service")
	cfg := config.NewConfig()
	ctx := context.Background()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	if err := a.Start(); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}
}
