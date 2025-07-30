package handler

import (
	"net/http"
	"recipe-processor/internal/models"
	"recipe-processor/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type RecipeHandler struct {
	service   service.RecipeServiceInterface
	validator *validator.Validate
	logger    *zap.Logger
}

func NewRecipeHandler(service service.RecipeServiceInterface, logger *zap.Logger) *RecipeHandler {
	return &RecipeHandler{
		service:   service,
		validator: validator.New(),
		logger:    logger,
	}
}

func (h *RecipeHandler) ProcessRecipe(c *gin.Context) {
	var req models.ProcessRecipeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ProcessRecipeResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Request validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ProcessRecipeResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
		return
	}

	_, err := h.service.ProcessRecipe(c.Request.Context(), req.Recipe)
	if err != nil {
		h.logger.Error("Failed to process recipe", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ProcessRecipeResponse{
			Success: false,
			Message: "Failed to process recipe: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ProcessRecipeResponse{
		Success: true,
		Message: "Recipe processed successfully",
	})
}
