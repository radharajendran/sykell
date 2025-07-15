package main

import (
	"os"
	"os/signal"
	"syscall"

	"sykell-backend/internal/router"
	"sykell-backend/internal/service"
	"sykell-backend/pkg/database"
	"sykell-backend/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Initialize database
	if err := database.Init(); err != nil {
		logger.Sugar().Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	// Start crawler job processor
	jobProcessor := service.NewCrawlerJobProcessor()
	jobProcessor.Start()
	defer jobProcessor.Stop()

	// Setup router
	app := router.Setup()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Sugar().Info("Gracefully shutting down...")
		jobProcessor.Stop()
		app.Shutdown()
	}()

	// Start server
	logger.Sugar().Info("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		logger.Sugar().Fatalf("Server error: %v", err)
	}
}
