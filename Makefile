.PHONY: proto proto-install proto-clean wire swagger build run test clean lint lint-frontend help

# Variables
PROTO_DIR=api/proto
PROTO_OUT_DIR=api/proto
GO_BIN=$(shell go env GOPATH)/bin

# Help command
help:
	@echo "Available commands:"
	@echo "  make proto-install     - Install protoc plugins"
	@echo "  make proto             - Generate Go code from proto files"
	@echo "  make proto-clean       - Clean generated proto files"
	@echo "  make wire              - Generate wire dependency injection code"
	@echo "  make swagger           - Generate Swagger documentation"
	@echo "  make build             - Build user service"
	@echo "  make run               - Run user service"
	@echo "  make test              - Run tests"
	@echo "  make lint              - Run Go linter"
	@echo "  make lint-fix          - Auto-fix Go formatting"
	@echo "  make lint-frontend     - Run frontend linter"
	@echo "  make lint-frontend-fix - Auto-fix frontend formatting"
	@echo "  make clean             - Clean build artifacts"

# Install protoc plugins
proto-install:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Protoc plugins installed successfully!"

# Generate Go code from proto files
proto:
	@echo "Generating Go code from proto files..."
	protoc --go_out=. --go_opt=module=github.com/ddd-micro \
		--go-grpc_out=. --go-grpc_opt=module=github.com/ddd-micro \
		--proto_path=. \
		$(PROTO_DIR)/**/*.proto
	@echo "Proto generation completed!"

# Clean generated proto files
proto-clean:
	@echo "Cleaning generated proto files..."
	find $(PROTO_OUT_DIR) -name "*.pb.go" -type f -delete
	@echo "Proto files cleaned!"

# Generate wire dependency injection code
wire:
	@echo "Generating wire code..."
	cd cmd/user && $(GO_BIN)/wire
	cd cmd/product && $(GO_BIN)/wire
	@echo "Wire generation completed!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/user/main.go -o cmd/user/docs
	@echo "Swagger documentation generated successfully!"

# Build services
build:
	@echo "Building services..."
	go build -o bin/user-service ./cmd/user
	go build -o bin/product-service ./cmd/product
	@echo "Build completed!"

# Run services
run-user:
	@echo "Running user service..."
	go run ./cmd/user/main.go

run-product:
	@echo "Running product service..."
	go run ./cmd/product/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run Go linter
lint:
	@echo "Running Go linter..."
	$(GO_BIN)/golangci-lint run ./...

# Auto-fix Go formatting
lint-fix:
	@echo "Auto-fixing Go formatting..."
	gofmt -w .
	$(GO_BIN)/goimports -w .
	@echo "✅ Go formatting fixed!"

# Run frontend linter
lint-frontend:
	@echo "Running frontend linter..."
	cd client && npm run lint
	cd client && npm run type-check
	cd client && npm run format:check

# Auto-fix frontend formatting
lint-frontend-fix:
	@echo "Auto-fixing frontend formatting..."
	cd client && npm run lint:fix
	cd client && npm run format
	@echo "✅ Frontend formatting fixed!"

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd client && npm install
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_BIN) v1.54.2

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf client/.next/
	rm -rf client/node_modules/
	@echo "Clean completed!"

# Full build (includes proto, wire, swagger generation and building)
full-build: proto wire swagger build
	@echo "Full build completed!"

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# Development commands
dev:
	@echo "Starting development environment..."
	docker-compose up -d user-db
	@echo "Database started. Run 'make run' to start the service."