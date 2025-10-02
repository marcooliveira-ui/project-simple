.PHONY: help run build test clean swagger docker-up docker-down migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run cmd/api/main.go

build: ## Build the application
	go build -o bin/api cmd/api/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf docs/
	rm -f coverage.out coverage.html

swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go -o docs

install-deps: ## Install dependencies
	go mod download
	go mod tidy

install-tools: ## Install development tools
	go install github.com/swaggo/swag/cmd/swag@latest

docker-up: ## Start PostgreSQL with docker-compose
	docker-compose up -d

docker-down: ## Stop docker-compose services
	docker-compose down

docker-logs: ## View docker logs
	docker-compose logs -f

migrate: ## Run database migrations (automatic on app start)
	@echo "Migrations run automatically when starting the application"

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

dev: swagger run ## Generate swagger and run the application
