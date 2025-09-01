.PHONY: dev up down build clean test logs build-lambda deploy-lambda

# Development with hot reload
dev:
	docker-compose up --build

# Start services in background
up:
	docker-compose up -d --build

# Stop services
down:
	docker-compose down

# Build production image
build:
	docker build --target production -t users_api:latest .

# Build server binary
build-server:
	go build -o bin/server cmd/server/main.go

# Clean up volumes and images
clean:
	docker-compose down -v
	docker system prune -f

# Clean Lambda artifacts
clean-lambda:
	rm -f bootstrap lambda.zip

# Run tests
test:
	go test ./...

# View logs
logs:
	docker-compose logs -f

# View app logs only
logs-app:
	docker-compose logs -f app

# View db logs only
logs-db:
	docker-compose logs -f db

# Build Lambda deployment package
build-lambda:
	@echo "Building Lambda deployment package..."
	@rm -f bootstrap lambda.zip
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/lambda/main.go
	zip lambda.zip bootstrap
	@echo "Lambda package created: lambda.zip"
	@ls -lh lambda.zip

# Deploy to Lambda (requires AWS CLI and function to exist)
deploy-lambda:
	./scripts/deploy-lambda.sh

# Full Lambda build and deploy
lambda: clean-lambda build-lambda deploy-lambda