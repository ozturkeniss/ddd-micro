.PHONY: proto proto-install proto-clean wire swagger build run test clean help

# Variables
PROTO_DIR=api/proto
PROTO_OUT_DIR=api/proto
GO_BIN=$(shell go env GOPATH)/bin

# Help command
help:
	@echo "Available commands:"
	@echo "  make proto-install  - Install protoc plugins"
	@echo "  make proto          - Generate Go code from proto files"
	@echo "  make proto-clean    - Clean generated proto files"
	@echo "  make wire           - Generate wire dependency injection code"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make build          - Build user service"
	@echo "  make run            - Run user service"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"

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
		$(PROTO_DIR)/user/user.proto
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
	@echo "Wire generation completed!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/user/main.go -o cmd/user/docs
	@echo "Swagger documentation generated successfully!"

# Build user service
build:
	@echo "Building user service..."
	go build -o bin/user-service ./cmd/user
	@echo "Build completed!"

# Run user service
run:
	@echo "Running user service..."
	go run ./cmd/user/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f cmd/user/wire_gen.go
	@echo "Clean completed!"

# Build docker image
docker-build:
	@echo "Building Docker image..."
	docker build -f dockerfiles/user.dockerfile -t user-service:latest .
	@echo "Docker build completed!"

# Run docker compose
docker-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

# Stop docker compose
docker-down:
	@echo "Stopping services..."
	docker-compose down

# Install all dependencies
install-deps: proto-install
	@echo "Installing Go dependencies..."
	go mod download
	go install github.com/google/wire/cmd/wire@latest
	@echo "All dependencies installed!"

# Full build (proto + wire + build)
full-build: proto wire build
	@echo "Full build completed!"

