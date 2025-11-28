# Makefile for Recipe Processor API

.PHONY: help build run test clean docker-build docker-run docker-stop lint fmt vet

# Variables
APP_NAME=recipe-processor
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_CONTAINER=$(APP_NAME)-container
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Colors for output
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}'

build: ## Build the application binary
	@echo "Building application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/api cmd/api/main.go
	@echo "${GREEN}Build complete: bin/api${NC}"

run: ## Run the application locally
	@echo "Running application..."
	go run cmd/api/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "${GREEN}Tests complete${NC}"

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report: coverage.html${NC}"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "${GREEN}Clean complete${NC}"

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "${GREEN}Format complete${NC}"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "${GREEN}Vet complete${NC}"

lint: ## Run golangci-lint (requires golangci-lint installed)
	@echo "Running linter..."
	golangci-lint run ./...
	@echo "${GREEN}Lint complete${NC}"

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	go mod tidy
	@echo "${GREEN}Tidy complete${NC}"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "${GREEN}Docker image built: $(DOCKER_IMAGE)${NC}"

docker-run: docker-build ## Build and run Docker container
	@echo "Starting Docker container..."
	docker run -d \
		--name $(DOCKER_CONTAINER) \
		-p 8080:8080 \
		-e ENV=production \
		-e PORT=8080 \
		$(DOCKER_IMAGE)
	@echo "${GREEN}Container started: $(DOCKER_CONTAINER)${NC}"
	@echo "API available at: http://localhost:8080"

docker-stop: ## Stop and remove Docker container
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_CONTAINER) 2>/dev/null || true
	docker rm $(DOCKER_CONTAINER) 2>/dev/null || true
	@echo "${GREEN}Container stopped${NC}"

docker-logs: ## Show Docker container logs
	docker logs -f $(DOCKER_CONTAINER)

docker-shell: ## Open shell in running container
	docker exec -it $(DOCKER_CONTAINER) sh

dev: ## Run in development mode with auto-reload (requires air)
	@echo "Starting development server with air..."
	air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "${GREEN}Tools installed${NC}"

.DEFAULT_GOAL := help