#=======================================================
# QKO BETA API - Development Operations
# Version: 1.0.0
#
# Required:
# - Docker 20.10.0+
# - Docker Compose 2.0.0+
#=======================================================

# Declare all phony targets
.PHONY: dev test down health health-detail logs logs-postgres logs-redis logs-mailhog troubleshoot-postgres troubleshoot-redis help clean init-env prod

# Environment type
ENVS := development test production
CURRENT_ENV ?= development

# Colors for output
BLUE := \033[34m
GREEN := \033[32m
RED := \033[31m
RESET := \033[0m

#==========================================
# Automates initial setup
#==========================================
# Add after your .PHONY declaration
init-env:
	@if [ ! -f .env.dev ]; then \
		cp .env.example .env.dev && \
		echo "$(GREEN)Created .env.dev from example$(RESET)"; \
	fi
	@if [ ! -f .env.test ]; then \
		cp .env.example .env.test && \
		echo "$(GREEN)Created .env.test from example$(RESET)"; \
	fi
	@echo "$(BLUE)Remember to update passwords in your .env files$(RESET)"


#==========================================
# Checks for docker and environment files
#==========================================

# Check required tools
check-docker:
	@if ! command -v docker > /dev/null 2>&1; then \
		echo "$(RED)Error: docker is not installed$(RESET)"; \
		exit 1; \
	fi
	@if ! docker compose version > /dev/null 2>&1; then \
		echo "$(RED)Error: Docker Compose is not installed$(RESET)"; \
		exit 1; \
	fi

# Check environment files
check-env-files:
	@if [ ! -f .env.dev ]; then \
		echo "$(RED)Error: .env.dev not found. Run 'make init-env' first$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f .env.test ]; then \
		echo "$(RED)Error: .env.test not found. Run 'make init-env' first$(RESET)"; \
		exit 1; \
	fi
	@if [ "$(CURRENT_ENV)" = "production" ] && [ ! -f .env.prod ]; then \
		echo "$(RED)Error: .env.prod not found$(RESET)"; \
		exit 1; \
	fi

#==========================================
# Environment Management
#==========================================

# Start development environment
dev: CURRENT_ENV=development
dev: check-docker check-env-files validate-env validate-env-type
	docker compose --env-file .env.dev up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@make health
	@echo "Launching Golang API..."
	@make run-api

# Start test environment
test: CURRENT_ENV=test
test: check-docker check-env-files validate-env validate-env-type
	docker compose --env-file .env.test up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@make health

# Start production environment
prod: CURRENT_ENV=production
prod: check-docker check-env-files validate-env validate-env-type
	@echo "$(RED)You are about to start the PRODUCTION environment$(RESET)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ ! $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Aborted$(RESET)"; \
		exit 1; \
	fi
	docker compose --env-file .env.prod up -d

# Shut down environment and remove volumes
down:
	docker compose --env-file .env.dev down -v

#==========================================
# Cleanup
#==========================================

# Remove all containers, volumes, and docker artifacts
clean:
	@echo "Cleaning up all resources..."
	docker compose --env-file .env.dev down -v
	docker system prune -f
	@echo "$(GREEN)Cleanup complete$(RESET)"

#==========================================
# Health Checks
#==========================================

# Quick health status check
health:
	docker compose ps

# Detailed health status with logs
health-detail:
	@echo "=== Container Status ==="
	docker compose ps
	@echo "\n=== Health Check Logs ==="
	docker compose ps | grep -q "healthy" || (echo "⚠️  Unhealthy services detected" && exit 1)
	@echo "✅ All services healthy"

#==========================================
# Logging
#==========================================

# View all service logs
logs:
	docker compose logs

# View Postgres logs
logs-postgres:
	docker compose logs postgres

# View Redis logs
logs-redis:
	docker compose logs redis

# View Mailhog logs
logs-mailhog:
	docker compose logs mailhog

#==========================================
# Troubleshooting
#==========================================

# Validate environment
validate-env:
	@if [ "$(CURRENT_ENV)" = "production" ]; then \
		if grep -qE '^POSTGRES_PASSWORD=CHANGE_ME' .env.prod; then \
			echo "$(RED)Error: Default production passwords detected. Please change them!$(RESET)"; \
			exit 1; \
		fi \
	fi

# Validate environment type
validate-env-type:
	@if ! echo "$(ENVS)" | grep -w "$(CURRENT_ENV)" > /dev/null; then \
		echo "$(RED)Error: CURRENT_ENV must be one of: $(ENVS)$(RESET)"; \
		exit 1; \
	fi
	@if [ "$(CURRENT_ENV)" != "production" ] && [ -f .env.prod ]; then \
		echo "$(RED)Error: .env.prod found in non-production environment$(RESET)"; \
		echo "$(RED)Please remove or rename .env.prod for development/test$(RESET)"; \
		exit 1; \
	fi

# Detailed Postgres troubleshooting
troubleshoot-postgres:
	@echo "=== Postgres Status ==="
	docker compose ps postgres
	@echo "\n=== Postgres Logs ==="
	docker compose logs --tail=50 postgres
	@echo "\n=== Postgres Environment ==="
	docker compose exec postgres env

# Detailed Redis troubleshooting
troubleshoot-redis:
	@echo "=== Redis Status ==="
	docker compose ps redis
	@echo "\n=== Redis Logs ==="
	docker compose logs --tail=50 redis
	@echo "\n=== Redis Ping Test ==="
	docker compose exec redis redis-cli ping

#==========================================
# Run API
#==========================================
# The run-api target starts the Golang API.
# It assumes that the main entrypoint of your app is located at ./cmd/api/main.go.
# Make sure your Go environment is properly configured to run the API.
run-api:
	@echo "Updating backend/.env from .env.dev..."
	@cp .env.dev backend/.env
	@echo "Starting Golang API..."
	cd backend && go run cmd/api/main.go


#==========================================
# Help
#==========================================

# Show this help
help:
	@echo "$(BLUE)Available commands:$(RESET)"
	@echo "$(BLUE)Environment:$(RESET)"
	@echo "  make init-env         - Initialize environment files from templates"
	@echo "  make test             - Start test environment"
	@echo "  make dev              - Start development environment"
	@echo "  make prod             - Start production environment"
	@echo "  make down             - Shut down environment"
	@echo "  make clean            - Remove all containers, volumes, and docker artifacts"
	@echo "\nHealth Checks:"
	@echo "  make health           - Quick health status"
	@echo "  make health-detail    - Detailed health status"
	@echo "\nLogging:"
	@echo "  make logs             - View all logs"
	@echo "  make logs-postgres    - View Postgres logs"
	@echo "  make logs-redis       - View Redis logs"
	@echo "  make logs-mailhog     - View Mailhog logs"
	@echo "\nTroubleshooting:"
	@echo "  make troubleshoot-postgres  - Troubleshoot Postgres"
	@echo "  make troubleshoot-redis     - Troubleshoot Redis"