package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"recipe-processor/internal/application/recipe"
	"recipe-processor/internal/domain"
	"recipe-processor/internal/infrastructure/http/handlers"
	"recipe-processor/internal/shared/logger"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockRecipeSubmitter is a mock implementation of RecipeSubmitter
type mockRecipeSubmitter struct {
	executeFunc func(ctx context.Context, cmd recipe.SubmitRecipeCommand) (*recipe.SubmitRecipeResult, error)
}

func (m *mockRecipeSubmitter) Execute(ctx context.Context, cmd recipe.SubmitRecipeCommand) (*recipe.SubmitRecipeResult, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, cmd)
	}

	// Default: simulate real service behavior (validate then succeed)
	if _, err := domain.NewRecipeText(cmd.RecipeText); err != nil {
		return nil, err
	}

	return &recipe.SubmitRecipeResult{RecipeID: "test-id"}, nil
}

func setupTestRouter(handler *handlers.RecipeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/recipes", handler.SubmitRecipe)
	return router
}

func TestRecipeHandler_SubmitRecipe_Success(t *testing.T) {
	// Arrange
	mockService := &mockRecipeSubmitter{
		executeFunc: func(ctx context.Context, cmd recipe.SubmitRecipeCommand) (*recipe.SubmitRecipeResult, error) {
			return &recipe.SubmitRecipeResult{RecipeID: "recipe-123"}, nil
		},
	}

	handler := handlers.NewRecipeHandler(logger.NewNoopLogger(), mockService)
	router := setupTestRouter(handler)

	reqBody := handlers.SubmitRecipeRequest{
		RecipeText: "Valid recipe text with ingredients and instructions",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, w.Code)
	}

	var response handlers.SubmitRecipeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.RecipeID != "recipe-123" {
		t.Errorf("Expected recipe_id 'recipe-123', got '%s'", response.RecipeID)
	}

	if response.Message != "Recipe submitted for processing" {
		t.Errorf("Expected message 'Recipe submitted for processing', got '%s'", response.Message)
	}
}

func TestRecipeHandler_SubmitRecipe_EmptyText(t *testing.T) {
	// Arrange
	mockService := &mockRecipeSubmitter{
		executeFunc: func(ctx context.Context, cmd recipe.SubmitRecipeCommand) (*recipe.SubmitRecipeResult, error) {
			return nil, domain.ErrRecipeTextEmpty
		},
	}

	handler := handlers.NewRecipeHandler(logger.NewNoopLogger(), mockService)
	router := setupTestRouter(handler)

	reqBody := handlers.SubmitRecipeRequest{RecipeText: ""}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "EMPTY_TEXT" {
		t.Errorf("Expected code 'EMPTY_TEXT', got '%s'", response.Code)
	}
}

func TestRecipeHandler_SubmitRecipe_TextTooLong(t *testing.T) {
	// Arrange
	mockService := &mockRecipeSubmitter{
		executeFunc: func(ctx context.Context, cmd recipe.SubmitRecipeCommand) (*recipe.SubmitRecipeResult, error) {
			return nil, domain.ErrRecipeTextTooLong
		},
	}

	handler := handlers.NewRecipeHandler(logger.NewNoopLogger(), mockService)
	router := setupTestRouter(handler)

	longText := strings.Repeat("a", 10001)
	reqBody := handlers.SubmitRecipeRequest{RecipeText: longText}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "TEXT_TOO_LONG" {
		t.Errorf("Expected code 'TEXT_TOO_LONG', got '%s'", response.Code)
	}
}

func TestRecipeHandler_SubmitRecipe_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := &mockRecipeSubmitter{}
	handler := handlers.NewRecipeHandler(logger.NewNoopLogger(), mockService)
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "INVALID_JSON" {
		t.Errorf("Expected code 'INVALID_JSON', got '%s'", response.Code)
	}
}

func TestRecipeHandler_SubmitRecipe_MissingRecipeText(t *testing.T) {
	// Arrange
	mockService := &mockRecipeSubmitter{}
	handler := handlers.NewRecipeHandler(logger.NewNoopLogger(), mockService)
	router := setupTestRouter(handler)

	reqBody := map[string]string{} // Missing recipe_text field
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// DEBUG: Print what we got back
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "EMPTY_TEXT" {
		t.Errorf("Expected code 'EMPTY_TEXT', got '%s'", response.Code)
	}
}

func TestRecipeHandler_SubmitRecipe_ServiceError(t *testing.T) {
	// Arrange
	mockService := &mockRecipeSubmitter{
		executeFunc: func(ctx context.Context, cmd recipe.SubmitRecipeCommand) (*recipe.SubmitRecipeResult, error) {
			return nil, errors.New("unexpected service error")
		},
	}

	handler := handlers.NewRecipeHandler(logger.NewNoopLogger(), mockService)
	router := setupTestRouter(handler)

	reqBody := handlers.SubmitRecipeRequest{RecipeText: "Valid recipe"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "INTERNAL_ERROR" {
		t.Errorf("Expected code 'INTERNAL_ERROR', got '%s'", response.Code)
	}
}
