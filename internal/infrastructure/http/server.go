package http

import (
	"context"
	"net/http"
	"recipe-processor/internal/application/recipe"
	"recipe-processor/internal/config"
	"recipe-processor/internal/infrastructure/http/handlers"
	"recipe-processor/internal/shared/events"
	"recipe-processor/internal/shared/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config   *config.Config
	logger   logger.Logger
	eventBus events.EventBus
	srv      *http.Server
}

func NewServer(cfg *config.Config, logger logger.Logger, eventBus events.EventBus) *Server {
	return &Server{
		config:   cfg,
		logger:   logger,
		eventBus: eventBus,
	}
}

func (s *Server) Start() error {
	router := s.setupRouter()

	s.srv = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(LoggingMiddleware(s.logger))
	router.Use(CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	v1 := router.Group("/api/v1")
	{
		// Recipe routes
		submitService := recipe.NewSubmitRecipeService(s.eventBus, s.logger)
		recipeHandler := handlers.NewRecipeHandler(s.logger, submitService)
		v1.POST("/recipes", recipeHandler.SubmitRecipe)
	}

	return router
}
