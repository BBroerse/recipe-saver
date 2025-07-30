package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"recipe-processor/internal/models"
	"time"

	"go.uber.org/zap"
)

type OllamaRepositoryInterface interface {
	ProcessRecipe(ctx context.Context, prompt string) (string, error)
}

type OllamaRepository struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewOllamaRepository(baseURL string, timeout time.Duration, logger *zap.Logger) *OllamaRepository {
	return &OllamaRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (r *OllamaRepository) ProcessRecipe(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", r.baseURL)

	reqBody := models.OllamaRequest{
		Model:  "mistral",
		Stream: false,
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	r.logger.Info("Sending request to Ollama", zap.String("url", url))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp models.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}
