.PHONY: dev up down build clean test logs

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

# Clean up volumes and images
clean:
	docker-compose down -v
	docker system prune -f

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