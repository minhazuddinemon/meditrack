.PHONY: dev build templ tailwind sqlc clean

# Default target runs the full dev environment
dev:
	@echo "Starting development mode..."
	@air

# Generate sqlc Go code from SQL files
sqlc:
	@echo "Generating sqlc code..."
	@sqlc generate

# Generate Templ components
templ:
	@echo "Generating Templ components..."
	@templ generate

# Watch and compile Tailwind CSS in standalone watch mode (optional backup)
tailwind-watch:
	@npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Build production binary
build: sqlc templ
	@echo "Building production CSS..."
	@npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
	@echo "Building Go application..."
	@go build -o bin/server ./cmd/server/main.go

# Clean temporary files
clean:
	@rm -rf tmp bin
