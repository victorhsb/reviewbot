package main

import (
	"context"

	"github.com/victorhsb/review-bot/backend/http/routes"
	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/service/processor"
	"github.com/victorhsb/review-bot/backend/service/workflows"
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

	messageProcessor := processor.NewLocalMessageProcessor(ctx)

	messageService := service.NewMessageService(repo, messageProcessor)
	productService := service.NewProductService(repo, service.ProductServiceConfig{DefaultListingLimit: 100})

	reviewWorkflow := workflows.NewReviewWorkflow(ctx, productService, messageService)
	cliWorkflow := workflows.NewCLIWorkflow(reviewWorkflow, messageService)

	messageProcessor.Use(reviewWorkflow)
	messageProcessor.Use(cliWorkflow)

	router := routes.NewRouter()
	routes.RegisterMessageRoutes(router, messageService)
	routes.RegisterProductRoutes(router, productService)

	if err := router.Run(cfg.Port); err != nil {
		panic(err)
	}
}
