package domain_test

import (
	"recipe-processor/internal/domain"

	"strings"
	"testing"
)

func TestNewRecipeText_Empty(t *testing.T) {
	_, err := domain.NewRecipeText("   ")
	if err != domain.ErrRecipeTextEmpty {
		t.Fatalf("ex[ected ErrRecipeTextEmpty, got %v", err)
	}
}

func TestNewRecipeText_TooLong(t *testing.T) {
	long := strings.Repeat("a", domain.MaxRecipeTextLength+1)
	_, err := domain.NewRecipeText(long)
	if err != domain.ErrRecipeTextTooLong {
		t.Fatalf("expected ErrRecipeTextTooLong, got %v", err)
	}
}

func TestNewRecipeText_Valid(t *testing.T) {
	raw := "  Pancakes with syrup  "
	rt, err := domain.NewRecipeText(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := "Pancakes with syrup"
	if got := rt.Value(); got != want {
		t.Fatalf("Value() = %q, want %q", got, want)
	}

	if s := rt.String(); s != want {
		t.Fatalf("String() = %q, want %q", s, want)
	}
}
