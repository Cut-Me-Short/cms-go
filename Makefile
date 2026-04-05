# CutMeShort Go SDK - Build and Development Makefile

.PHONY: help test lint fmt vet build clean coverage docs examples

help:
	@echo "CutMeShort Go SDK - Development Commands"
	@echo "========================================="
	@echo "make test       - Run all tests"
	@echo "make test-v     - Run tests with verbose output"
	@echo "make coverage   - Run tests with coverage report"
	@echo "make lint       - Run golangci-lint"
	@echo "make fmt        - Format code with gofmt"
	@echo "make vet        - Run go vet"
	@echo "make build      - Build examples"
	@echo "make clean      - Clean build artifacts"
	@echo "make deps       - Download and verify dependencies"
	@echo "make security   - Run security checks"

# Testing
test:
	go test -v -race ./...

test-v:
	go test -v -race -count=1 ./...

coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Linting and Formatting
lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	go mod tidy

vet:
	go vet ./...

# Security
security:
	gosec ./...
	go run github.com/google/osv-scanner/cmd/osv-scanner@latest -r .

# Building
build: examples

examples:
	@echo "Building examples..."
	go build -o bin/track_lead examples/track_lead.go
	go build -o bin/track_sale examples/track_sale.go
	go build -o bin/basic_config examples/basic_config.go
	@echo "Built: bin/track_lead bin/track_sale bin/basic_config"

# Dependencies
deps:
	go mod download
	go mod verify
	go mod tidy

# Cleanup
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "Cleaned build artifacts"

# All checks (CI/CD equivalent)
all: fmt vet lint test coverage security
	@echo "All checks completed successfully"
