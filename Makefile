APP_NAME := recipe-saver
CMD_PATH := ./cmd/recipe-saver

.PHONY: help lint test build run clean

help:
	@echo "Available commands:"
	@echo " make install-tools Install all tools"
	@echo " make lint          Run linters"
	@echo " make test          Run unit tests"
	@echo " make build         Compile the Go binary"
	@echo " make run           Run the application"
	@echo " make clean         Remove build artifacts"

install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/conventionalcommit/commitlint@latest

lint:
	golangci-lint run

test:
	go test -v ./...

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) .

run:
	go run main.go

clean:
	rm -rf bin/