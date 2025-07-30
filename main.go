package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"recipe-processor/internal/config"
	"recipe-processor/internal/handler"
	"recipe-processor/internal/repository"
	"recipe-processor/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	defer logger.Sync() // nolint:errcheck

	cfg := config.Load()
	deps := initializeDependencies(cfg, logger)
	router := setupRouter(deps)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

type Dependencies struct {
	RecipeHandler *handler.RecipeHandler
	Logger        *zap.Logger
}

func initializeDependencies(cfg *config.Config, logger *zap.Logger) *Dependencies {
	ollamaRepo := repository.NewOllamaRepository(cfg.OllamaBaseUrl, cfg.RequestTimeout, logger)

	recipeService := service.NewRecipeService(ollamaRepo, logger)

	recipeHandler := handler.NewRecipeHandler(recipeService, logger)

	return &Dependencies{
		RecipeHandler: recipeHandler,
		Logger:        logger,
	}
}

func setupRouter(deps *Dependencies) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API routes
	router.POST("/api/v1/recipes/process", deps.RecipeHandler.ProcessRecipe)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
