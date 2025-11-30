package domain_test

import (
	"recipe-processor/internal/domain"
	"testing"
	"time"
)

func TestNewRecipeSubmitted_EventTypeAndTimestamp(t *testing.T) {
	e := domain.NewRecipeSubmitted("recipe-123", "some text")

	if et := e.EventType(); et != domain.EventTypeRecipeSubmitted {
		t.Fatalf("EventType() = %q, want %q", et, domain.EventTypeRecipeSubmitted)
	}

	if e.RecipeID != "recipe-123" {
		t.Fatalf("RecipeID = %q, want %q", e.RecipeID, "recipe-123")
	}

	// OccurredAt should be recent
	if time.Since(e.OccurredAt()) > 5*time.Second {
		t.Fatalf("OccurredAt is too far in the past: %v", e.OccurredAt())
	}
}
