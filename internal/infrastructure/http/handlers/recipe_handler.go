package handlers

import (
	"errors"
	"net/http"
	"recipe-processor/internal/application/recipe"
	"recipe-processor/internal/domain"
	"recipe-processor/internal/shared/logger"

	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	logger        logger.Logger
	submitService *recipe.SubmitRecipeService
}

func NewRecipeHandler(logger logger.Logger, submitService *recipe.SubmitRecipeService) *RecipeHandler {
	return &RecipeHandler{
		logger:        logger,
		submitService: submitService,
	}
}

type SubmitRecipeRequest struct {
	RecipeText string `json:"recipe_text" binding:"required"`
}

type SubmitRecipeResponse struct {
	RecipeID string `json:"recipe_id"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// SubmitRecipe handles POST /api/v1/recipes
func (h *RecipeHandler) SubmitRecipe(c *gin.Context) {
	var req SubmitRecipeRequest

	// Parse and validate HTTP request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}

	// Execute business logic via application service
	cmd := recipe.SubmitRecipeCommand{
		RecipeText: req.RecipeText,
	}

	result, err := h.submitService.Execute(c.Request.Context(), cmd)
	if err != nil {
		// Map domain errors to HTTP responses
		statusCode, errorResp := h.mapErrorToResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	// Return success response
	c.JSON(http.StatusAccepted, SubmitRecipeResponse{
		RecipeID: result.RecipeID,
		Message:  "Recipe submitted for processing",
	})
}

func (h *RecipeHandler) mapErrorToResponse(err error) (int, ErrorResponse) {
	// Check for domain validation errors
	if errors.Is(err, domain.ErrRecipeTextEmpty) {
		return http.StatusBadRequest, ErrorResponse{
			Error:   "Recipe text validation failed",
			Code:    "EMPTY_TEXT",
			Details: err.Error(),
		}
	}

	if errors.Is(err, domain.ErrRecipeTextTooLong) {
		return http.StatusBadRequest, ErrorResponse{
			Error:   "Recipe text validation failed",
			Code:    "TEXT_TOO_LONG",
			Details: err.Error(),
		}
	}

	// Default to internal server error
	h.logger.Error("Unexpected error in recipe submission", logger.Error(err))
	return http.StatusInternalServerError, ErrorResponse{
		Error: "Failed to process recipe",
		Code:  "INTERNAL_ERROR",
	}
}
