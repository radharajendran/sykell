.PHONY: build run test clean setup-db

# Build the application
build:
	go build -o bin/sykell-api main.go

# Run the application
run:
	go run main.go

# Run with air for hot reloading (requires: go install github.com/cosmtrek/air@latest)
dev:
	air

# Test the application
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Setup database (requires MySQL to be running)
setup-db:
	mysql -u root -p < scripts/setup_db.sql

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires: go install honnef.co/go/tools/cmd/staticcheck@latest)
lint:
	staticcheck ./...

# Create binary directory
bin:
	mkdir -p bin

# Full build with binary directory
full-build: bin build

# Test API (requires the server to be running)
test-api:
	./scripts/test_api.sh

# Test Crawler API (requires the server to be running)
test-crawler:
	./scripts/test_crawler_api.sh
