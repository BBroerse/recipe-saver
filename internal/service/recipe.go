package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"recipe-processor/internal/models"
	"recipe-processor/internal/repository"

	"go.uber.org/zap"
)

type RecipeServiceInterface interface {
	ProcessRecipe(ctx context.Context, recipeData string) (*models.ProcessedRecipe, error)
}

type RecipeService struct {
	ollamaRepo repository.OllamaRepositoryInterface
	logger     *zap.Logger
}

func NewRecipeService(
	ollamaRepo repository.OllamaRepositoryInterface,
	logger *zap.Logger,
) *RecipeService {
	return &RecipeService{
		ollamaRepo: ollamaRepo,
		logger:     logger,
	}
}

func (s *RecipeService) ProcessRecipe(ctx context.Context, recipeData string) (*models.ProcessedRecipe, error) {
	s.logger.Info("Processing recipe", zap.Int("length", len(recipeData)))

	processedRecipe, err := s.processWithOllama(ctx, recipeData)
	if err != nil {
		s.logger.Error("Failed to process with Ollama", zap.Error(err))
		return nil, fmt.Errorf("ollama processing failed: %w", err)
	}

	s.logger.Info("Recipe processed successfully", zap.String("title", processedRecipe.Title))
	return processedRecipe, nil
}

func (s *RecipeService) processWithOllama(ctx context.Context, recipeData string) (*models.ProcessedRecipe, error) {
	prompt := s.buildOllamaPrompt(recipeData)

	response, err := s.ollamaRepo.ProcessRecipe(ctx, prompt)
	if err != nil {
		return nil, err
	}

	processedRecipe, err := s.parseOllamaResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ollama response: %w", err)
	}

	return processedRecipe, nil
}

func (s *RecipeService) buildOllamaPrompt(recipeData string) string {
	return fmt.Sprintf(`
		Please analyze the following recipe and extract structured information.
		All texts returned to me must be in dutch!

		{
			\"title\": \"Recipe title\",
		 	\"ingredients\": [\"ingredient 1\", \"ingredient 2\"],
			\"instructions\": [\"step 1\", \"step 2\"],
			\"total_time_minutes\": "15",
			\"servings\": "4",
			\"course_type\": \"main\"
		} Recipe data

		Recipe data: %s`, recipeData)
}

func (s *RecipeService) parseOllamaResponse(response string) (*models.ProcessedRecipe, error) {
	// Find JSON in response (in case there's extra text)
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[start : end+1]
	s.logger.Info("Ollama response: " + jsonStr)

	var parsed struct {
		Title        string   `json:"title"`
		Ingredients  []string `json:"ingredients"`
		Instructions []string `json:"instructions"`
		TotalTime    string   `json:"total_time"`
		Servings     string   `json:"servings"`
		CourseType   string   `json:"course_type"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert to internal model
	recipe := &models.ProcessedRecipe{
		Title:        parsed.Title,
		Ingredients:  parsed.Ingredients,
		Instructions: parsed.Instructions,
		Servings:     parsed.Servings,
		Totaltime:    parsed.TotalTime,
	}

	return recipe, nil
}
