package domain

import "time"

const (
	EventTypeRecipeSubmitted = "recipe.submitted"
)

type RecipeSubmitted struct {
	RecipeID   string
	RecipeText string
	occurredAt time.Time
}

// NewRecipeSubmitted creates a new RecipeSubmitted event
func NewRecipeSubmitted(recipeID, recipeText string) *RecipeSubmitted {
	return &RecipeSubmitted{
		RecipeID:   recipeID,
		RecipeText: recipeText,
		occurredAt: time.Now(),
	}
}

// EventType implements Event interface
func (e *RecipeSubmitted) EventType() string {
	return EventTypeRecipeSubmitted
}

// OccurredAt implements Event interface
func (e *RecipeSubmitted) OccurredAt() time.Time {
	return e.occurredAt
}
