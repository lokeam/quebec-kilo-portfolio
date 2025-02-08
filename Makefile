# QKO BETA API - Development Operations (Containerized Setup)
# Version: 2.0.0
#
# This Makefile leverages Docker Compose to build, run, and manage
# your containerized services:
#   - API (your Golang application)
#   - Redis
#   - Postgres
#   - Mailhog
#
# Best practices used:
#   - All services run in containers on a shared private network.
#   - The API container connects to Redis via Docker’s internal DNS (using the service name).
#   - Environment-specific configuration is loaded by using .env files.
#   - No sensitive service is exposed publicly (only the API is exposed with a published port).
#

.PHONY: init-env check-docker check-env-files dev test prod down clean health health-detail logs logs-postgres logs-redis logs-mailhog troubleshoot-postgres troubleshoot-redis help

# Define allowed environments and set current environment
ENVS := development test production
CURRENT_ENV ?= development

# Colors (for terminal output)
BLUE := \033[34m
GREEN := \033[32m
RED := \033[31m
RESET := \033[0m

# -----------------------------------------
# Environment Initialization
# -----------------------------------------
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

# -----------------------------------------
# Pre-check Tasks
# -----------------------------------------
check-docker:
	@if ! command -v docker > /dev/null 2>&1; then \
		echo "$(RED)Error: docker is not installed$(RESET)"; \
		exit 1; \
	fi
	@if ! docker compose version > /dev/null 2>&1; then \
		echo "$(RED)Error: Docker Compose is not installed$(RESET)"; \
		exit 1; \
	fi

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

# -----------------------------------------
# Environment Management via Docker Compose
# -----------------------------------------
dev: CURRENT_ENV=development
dev: check-docker check-env-files
	@echo "$(BLUE)Starting development environment...$(RESET)"
	docker compose --env-file .env.dev up --build -d
	@sleep 5
	@$(MAKE) health

test: CURRENT_ENV=test
test: check-docker check-env-files
	@echo "$(BLUE)Starting test environment...$(RESET)"
	docker compose --env-file .env.test up --build -d
	@sleep 5
	@$(MAKE) health

prod: CURRENT_ENV=production
prod: check-docker check-env-files
	@echo "$(RED)You are about to start the PRODUCTION environment$(RESET)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ ! $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Aborted$(RESET)"; \
		exit 1; \
	fi
	docker compose --env-file .env.prod up --build -d

down:
	docker compose --env-file .env.dev down -v

clean:
	@echo "$(BLUE)Cleaning up all resources...$(RESET)"
	docker compose --env-file .env.dev down -v
	docker system prune -f
	@echo "$(GREEN)Cleanup complete$(RESET)"

# -----------------------------------------
# Health Checks and Logs
# -----------------------------------------
health:
	docker compose ps

health-detail:
	@echo "=== Container Status ==="
	docker compose ps
	@echo "\n=== Health Check Logs ==="
	docker compose ps | grep -q "healthy" || (echo "⚠️ Unhealthy services detected" && exit 1)
	@echo "✅ All services healthy"

logs:
	docker compose logs

logs-postgres:
	docker compose logs postgres

logs-redis:
	docker compose logs redis

logs-mailhog:
	docker compose logs mailhog

# -----------------------------------------
# Troubleshooting
# -----------------------------------------
troubleshoot-postgres:
	@echo "=== Postgres Status ==="
	docker compose ps postgres
	@echo "\n=== Postgres Logs ==="
	docker compose logs --tail=50 postgres
	@echo "\n=== Postgres Environment ==="
	docker compose exec postgres env

troubleshoot-redis:
	@echo "=== Redis Status ==="
	docker compose ps redis
	@echo "\n=== Redis Logs ==="
	docker compose logs --tail=50 redis
	@echo "\n=== Redis Ping Test ==="
	docker compose exec redis redis-cli ping

# -----------------------------------------
# Help
# -----------------------------------------
help:
	@echo "$(BLUE)Available commands:$(RESET)"
	@echo "$(BLUE)Environment:$(RESET)"
	@echo "  make init-env         - Initialize environment files from templates"
	@echo "  make dev              - Start development environment (containerized)"
	@echo "  make test             - Start test environment (containerized)"
	@echo "  make prod             - Start production environment (containerized)"
	@echo "  make down             - Shut down environment and remove volumes"
	@echo "  make clean            - Remove all containers, volumes, and Docker artifacts"
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