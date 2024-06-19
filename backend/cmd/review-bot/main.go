package main

import (
	"context"

	"github.com/victorhsb/review-bot/backend/http/routes"
	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/storage/postgres"
)

func main() {
	ctx := context.Background()
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	repo, err := postgres.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err := repo.Migrate(); err != nil {
		panic(err)
	}

	messsageService := service.NewMessageService(repo)
	productService := service.NewProductService(repo, service.ProductServiceConfig{DefaultListingLimit: 100})

	router := routes.NewRouter()
	routes.RegisterMessageRoutes(router, messsageService)
	routes.RegisterProductRoutes(router, productService)

	if err := router.Run(cfg.Port); err != nil {
		panic(err)
	}
}
