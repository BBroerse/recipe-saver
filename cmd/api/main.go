package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"recipe-processor/internal/config"
	"recipe-processor/internal/infrastructure/http"
	"recipe-processor/internal/shared/events"
	"recipe-processor/internal/shared/logger"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	appLogger, err := logger.NewZapLogger(cfg.Environment)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	defer func() {
		if zapLogger, ok := appLogger.(*logger.ZapLogger); ok {
			_ = zapLogger.Sync()
		}
	}()

	appLogger.Info("Starting application",
		logger.String("environment", cfg.Environment),
		logger.String("port", cfg.Port),
	)

	// Initialize event bus (in-memory implementation)
	eventBus := events.NewMemoryEventBus(appLogger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start event bus
	if err := eventBus.Start(ctx); err != nil {
		appLogger.Fatal("Failed to start event bus", logger.Error(err))
	}

	server := http.NewServer(cfg, appLogger, eventBus)

	go func() {
		appLogger.Info("Starting server", logger.String("port", cfg.Port))

		if err := server.Start(); err != nil {
			appLogger.Fatal("Failed to start server", logger.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop server
	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Server shutdown error", logger.Error(err))
	}

	appLogger.Info("Server exited")
}
