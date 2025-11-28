package domain

import (
	"errors"
	"strings"
)

const (
	MaxRecipeTextLength = 10_000
)

var (
	ErrRecipeTextEmpty   = errors.New("recipe text cannot be empty")
	ErrRecipeTextTooLong = errors.New("recipe text exceeds maximum length")
)

type RecipeText struct {
	value string
}

func NewRecipeText(text string) (*RecipeText, error) {
	text = strings.TrimSpace(text)

	if text == "" {
		return nil, ErrRecipeTextEmpty
	}

	if len(text) > MaxRecipeTextLength {
		return nil, ErrRecipeTextTooLong
	}

	return &RecipeText{value: text}, nil
}

// Value returns the validated recipe text
func (rt *RecipeText) Value() string {
	return rt.value
}

// String implements Stringer interface
func (rt *RecipeText) String() string {
	return rt.value
}
