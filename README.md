# Users API - Development Environment

Go REST API with Clean Architecture + DDD + DI using Docker and Air for hot reload development.

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Make (optional, for convenience)

### Start Development Environment

```bash
# Using Make
make dev

# Or directly with docker-compose
docker-compose up --build
```

The API will be available at `http://localhost:8080`

### API Endpoints

- `GET /{tenant_id}/Users` - Search users with optional filters

Example:
```bash
# Get all users for tenant1
curl "http://localhost:8080/tenant1/Users"

# Search with filters
curl "http://localhost:8080/tenant1/Users?user_name=john&limit=10&offset=0"
```

### Development Commands

```bash
# Start in background
make up

# Stop services
make down

# View logs
make logs

# View only app logs
make logs-app

# Clean up everything
make clean

# Run tests
make test
```

### Hot Reload

The development setup uses [Air](https://github.com/cosmtrek/air) for automatic rebuilding and restarting when Go files change.

### Database

PostgreSQL runs in Docker with sample data for `tenant1` and `tenant2`. The database is initialized with test users on first run.

### Project Structure

```
src/
├── domain/             # Entities and repository interfaces
├── usecase/            # Business logic
├── interface/web/      # HTTP handlers and routing (gin)
├── infrastructure/     # Database implementation (PostgreSQL)
├── runtime/            # Dependency injection (fx)
└── main.go            # Application entry point
```

### Environment Variables

Set in `docker-compose.yml`:
- `DB_HOST=db`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=app`
- `DB_SSLMODE=disable`