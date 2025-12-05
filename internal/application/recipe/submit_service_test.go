package recipe_test

import (
	"context"
	"errors"
	"recipe-processor/internal/application/recipe"
	"recipe-processor/internal/domain"
	"recipe-processor/internal/shared/events"
	"recipe-processor/internal/shared/logger"
	"strings"
	"testing"
)

// mockEventBus is a mock implementation of EventBus
type mockEventBus struct {
	publishFunc   func(ctx context.Context, event events.Event) error
	publishCalled bool
	lastEvent     events.Event
}

func (m *mockEventBus) Publish(ctx context.Context, event events.Event) error {
	m.publishCalled = true
	m.lastEvent = event
	if m.publishFunc != nil {
		return m.publishFunc(ctx, event)
	}
	return nil
}

func (m *mockEventBus) Subscribe(eventType string, handler events.EventHandler) {}

func (m *mockEventBus) Start(ctx context.Context) error {
	return nil
}

func (m *mockEventBus) Stop() error {
	return nil
}

func TestSubmitRecipeService_Execute_Success(t *testing.T) {
	// Arrange
	mockBus := &mockEventBus{}
	log := logger.NewNoopLogger()
	service := recipe.NewSubmitRecipeService(mockBus, log)

	cmd := recipe.SubmitRecipeCommand{
		RecipeText: "Chocolate Chip Cookies\n\nIngredients:\n- 2 cups flour\n- 1 cup sugar\n\nInstructions:\n1. Mix ingredients\n2. Bake at 350F",
	}

	// Act
	result, err := service.Execute(context.Background(), cmd)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.RecipeID == "" {
		t.Error("Expected non-empty recipe ID")
	}

	if !mockBus.publishCalled {
		t.Error("Expected Publish to be called")
	}

	if mockBus.lastEvent == nil {
		t.Fatal("Expected event to be published")
	}

	// Verify event type
	if mockBus.lastEvent.EventType() != domain.EventTypeRecipeSubmitted {
		t.Errorf("Expected event type '%s', got '%s'",
			domain.EventTypeRecipeSubmitted,
			mockBus.lastEvent.EventType())
	}

	// Verify event data
	recipeEvent, ok := mockBus.lastEvent.(*domain.RecipeSubmitted)
	if !ok {
		t.Fatal("Expected RecipeSubmitted event")
	}

	if recipeEvent.RecipeID != result.RecipeID {
		t.Errorf("Expected event recipe ID '%s', got '%s'",
			result.RecipeID, recipeEvent.RecipeID)
	}

	if recipeEvent.RecipeText != cmd.RecipeText {
		t.Error("Expected event to contain original recipe text")
	}
}

func TestSubmitRecipeService_Execute_EmptyText(t *testing.T) {
	// Arrange
	mockBus := &mockEventBus{}
	log := logger.NewNoopLogger()
	service := recipe.NewSubmitRecipeService(mockBus, log)

	cmd := recipe.SubmitRecipeCommand{
		RecipeText: "",
	}

	// Act
	result, err := service.Execute(context.Background(), cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error for empty text, got nil")
	}

	if !errors.Is(err, domain.ErrRecipeTextEmpty) {
		t.Errorf("Expected error to wrap ErrRecipeTextEmpty, got: %v", err)
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	if mockBus.publishCalled {
		t.Error("Expected Publish not to be called on validation error")
	}
}

func TestSubmitRecipeService_Execute_TextTooLong(t *testing.T) {
	// Arrange
	mockBus := &mockEventBus{}
	log := logger.NewNoopLogger()
	service := recipe.NewSubmitRecipeService(mockBus, log)

	cmd := recipe.SubmitRecipeCommand{
		RecipeText: strings.Repeat("a", 10001), // Over 10,000 char limit
	}

	// Act
	result, err := service.Execute(context.Background(), cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error for text too long, got nil")
	}

	if !errors.Is(err, domain.ErrRecipeTextTooLong) {
		t.Errorf("Expected error to wrap ErrRecipeTextTooLong, got: %v", err)
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	if mockBus.publishCalled {
		t.Error("Expected Publish not to be called on validation error")
	}
}

func TestSubmitRecipeService_Execute_PublishError(t *testing.T) {
	// Arrange
	publishError := errors.New("event bus connection failed")
	mockBus := &mockEventBus{
		publishFunc: func(ctx context.Context, event events.Event) error {
			return publishError
		},
	}
	log := logger.NewNoopLogger()
	service := recipe.NewSubmitRecipeService(mockBus, log)

	cmd := recipe.SubmitRecipeCommand{
		RecipeText: "Valid recipe text",
	}

	// Act
	result, err := service.Execute(context.Background(), cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error when publish fails, got nil")
	}

	if !errors.Is(err, publishError) {
		t.Errorf("Expected error to wrap publish error, got: %v", err)
	}

	if result != nil {
		t.Error("Expected nil result when publish fails")
	}

	if !mockBus.publishCalled {
		t.Error("Expected Publish to be called")
	}
}

func TestSubmitRecipeService_Execute_ContextCancellation(t *testing.T) {
	// Arrange
	mockBus := &mockEventBus{
		publishFunc: func(ctx context.Context, event events.Event) error {
			// Simulate publish checking context
			return ctx.Err()
		},
	}
	log := logger.NewNoopLogger()
	service := recipe.NewSubmitRecipeService(mockBus, log)

	cmd := recipe.SubmitRecipeCommand{
		RecipeText: "Valid recipe text",
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Act
	result, err := service.Execute(ctx, cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error with cancelled context, got nil")
	}

	if result != nil {
		t.Error("Expected nil result with cancelled context")
	}
}

func TestSubmitRecipeService_Execute_RecipeIDIsUnique(t *testing.T) {
	// Arrange
	mockBus := &mockEventBus{}
	log := logger.NewNoopLogger()
	service := recipe.NewSubmitRecipeService(mockBus, log)

	cmd := recipe.SubmitRecipeCommand{
		RecipeText: "Recipe text",
	}

	// Act - Execute multiple times
	result1, err1 := service.Execute(context.Background(), cmd)
	result2, err2 := service.Execute(context.Background(), cmd)

	// Assert
	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}

	if result1.RecipeID == result2.RecipeID {
		t.Error("Expected unique recipe IDs for different submissions")
	}
}
