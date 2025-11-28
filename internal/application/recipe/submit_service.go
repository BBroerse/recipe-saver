package recipe

import (
	"context"
	"fmt"
	"recipe-processor/internal/domain"
	"recipe-processor/internal/shared/events"
	"recipe-processor/internal/shared/logger"

	"github.com/google/uuid"
)

// SubmitRecipeCommand represents the input for submitting a recipe
type SubmitRecipeCommand struct {
	RecipeText string
}

// SubmitRecipeResult represents the output of submitting a recipe
type SubmitRecipeResult struct {
	RecipeID string
}

// SubmitRecipeService handles the business logic for submitting recipes
type SubmitRecipeService struct {
	eventBus events.EventBus
	logger   logger.Logger
}

// NewSubmitRecipeService creates a new recipe submission service
func NewSubmitRecipeService(eventBus events.EventBus, log logger.Logger) *SubmitRecipeService {
	return &SubmitRecipeService{
		eventBus: eventBus,
		logger:   log,
	}
}

// Execute processes a recipe submission command
func (s *SubmitRecipeService) Execute(ctx context.Context, cmd SubmitRecipeCommand) (*SubmitRecipeResult, error) {
	// Validate recipe text using domain value object
	recipeText, err := domain.NewRecipeText(cmd.RecipeText)
	if err != nil {
		return nil, fmt.Errorf("recipe text validation failed: %w", err)
	}

	// Generate unique recipe ID
	recipeID := uuid.New().String()
	event := domain.NewRecipeSubmitted(recipeID, recipeText.Value())

	// Publish event
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish recipe submitted event",
			logger.String("recipe_id", recipeID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	s.logger.Info("Recipe submitted successfully",
		logger.String("recipe_id", recipeID),
		logger.Int("text_length", len(recipeText.Value())),
	)

	return &SubmitRecipeResult{
		RecipeID: recipeID,
	}, nil
}
