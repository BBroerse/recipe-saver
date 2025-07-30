package models

// Request/Response models
type ProcessRecipeRequest struct {
	Recipe string `json:"recipe" validate:"required,min=10,max=50000"`
}

type ProcessRecipeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Internal models
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type ProcessedRecipe struct {
	Title        string   `json:"title"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Totaltime    string   `json:"total_time"`
	Servings     string   `json:"servings"`
	CourseType   string   `json:"course_type"`
}
