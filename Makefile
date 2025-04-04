# QKO BETA API - Development Operations (Containerized Setup)
# Version: 2.7.0
#
# This Makefile leverages Docker Compose to build, run, and manage
# the following containerized services:
#   - Backend API (Golang application)
#   - Frontend (React application)
#   - Redis
#   - Postgres
#   - Mailhog
#   - Prometheus (Metrics collection)
#   - Grafana (Metrics visualization)
#   - Sentry (Backend Error tracking)
#
# Best practices used:
#   - All services run in containers on a shared private network.
#   - The API container connects to Redis via Docker's internal DNS (using the service name).
#   - Environment-specific configuration is loaded by using .env files.
#   - No sensitive service is exposed publicly (only the API is exposed with a published port).
#   - Metrics are collected and visualized for monitoring system health.
#

.PHONY: init-env check-docker check-env-files dev test prod down clean health health-detail logs logs-postgres logs-redis logs-mailhog logs-prometheus logs-grafana troubleshoot-postgres troubleshoot-redis troubleshoot-prometheus troubleshoot-grafana monitoring monitoring-down verify-sentry run-with-sentry test-sentry dev-with-sentry help

# Define allowed environments and set current environment
ENVS := development test production
CURRENT_ENV ?= development
REDIS_VOLUME ?= redis_data_updated

# Colors (for terminal output)
BLUE := \033[34m
GREEN := \033[32m
RED := \033[31m
RESET := \033[0m

# ---------------------------------------------------------------------------
# Reset: Forcefully tear down containers + remove persistent volumes
#
# This target stops all running containers and removes all associated volumes,
# including forcing the removal of the designated Redis volume. It also prunes
# any dangling volumes from the Docker system.
#
# Usage: make reset
#
# ⚠️ Warning: This is a destructive operation that removes persistent data.
# ---------------------------------------------------------------------------
reset:
	@echo "Stopping containers and removing volumes..."
	docker compose down -v
	@echo "Forcibly removing redis volume..."
	docker volume rm -f $(REDIS_VOLUME)
	@echo "Pruning dangling volumes (if any)..."
	docker volume prune -f

# ---------------------------------------------------------------------------
# Up: Build + start backend services (excluding frontend)
#
# This target rebuilds and starts the backend services along with necessary
# dependencies, including:
#   • Traefik (as the reverse proxy)
#   • Backend API (Golang application)
#   • Redis
#   • Postgres
#   • Mailhog
#
# The frontend service is purposefully excluded when using this target.
#
# Usage: make up
# ---------------------------------------------------------------------------
up:
	@echo "Rebuilding and starting backend services (excluding frontend)..."
	docker compose --env-file .env up --build -d traefik api redis postgres mailhog

# ---------------------------------------------------------------------------
# Restart: Fully reset the environment + bring up the backend services
#
# This composite target first calls the 'reset' target to clean up existing
# containers and volumes, then invokes the 'up' target to rebuild and start
# the backend services.
#
# Usage: make restart
#
# ⚠️ Warning: This process is destructive as it removes all running containers
# and volumes before restarting.
# ---------------------------------------------------------------------------
restart: reset up

clean:
	docker compose down

# ---------------------------------------------------------------------------
# Restart-All: Fully reset the environment + bring up ALL services (backend & frontend)
#
# This composite target calls the 'reset' target to clean up existing
# containers and volumes, then starts both backend and frontend services.
#
# Usage: make restart-all
#
# ⚠️ Warning: This process is destructive as it removes all running containers
# and volumes before restarting all services.
# ---------------------------------------------------------------------------
restart-all: CURRENT_ENV=development
restart-all: check-docker check-env-files
	@echo "$(BLUE)Fully resetting environment and starting all services (backend & frontend)...$(RESET)"
	$(MAKE) reset
	docker compose --env-file .env.dev up --build -d
	@sleep 5
	@$(MAKE) health

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

# Start development environment
dev: CURRENT_ENV=development
dev: check-docker check-env-files
	@echo "$(BLUE)Starting development environment...$(RESET)"
	docker compose --env-file .env.dev up --build -d
	@sleep 5
	@$(MAKE) health

# Start test environment
test: CURRENT_ENV=test
test: check-docker check-env-files
	@echo "$(BLUE)Starting test environment...$(RESET)"
	docker compose --env-file .env.test up --build -d
	@sleep 5
	@$(MAKE) health

# Start production environment
prod: CURRENT_ENV=production
prod: check-docker check-env-files
	@echo "$(RED)You are about to start the PRODUCTION environment$(RESET)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ ! $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Aborted$(RESET)"; \
		exit 1; \
	fi
	docker compose --env-file .env.prod up --build -d

# -----------------------------------------
# Start development environment without building the frontend
# -----------------------------------------
dev-backend: CURRENT_ENV=development
dev-backend: check-docker check-env-files
	@echo "$(BLUE)Starting development environment without frontend...$(RESET)"
	docker compose --env-file .env.dev up --build -d traefik api postgres redis mailhog
	@sleep 5
	@$(MAKE) health

# Shut down the development environment
down:
	docker compose --env-file .env.dev down -v

# Clean up all resources
clean:
	@echo "$(BLUE)Cleaning up all resources...$(RESET)"
	docker compose --env-file .env.dev down -v
	docker system prune -f
	@echo "$(GREEN)Cleanup complete$(RESET)"

# ---------------------------------------------------------------------------
# Full-Stack: Start all services including monitoring
#
# This target starts all services including the monitoring stack.
#
# Usage: make full-stack
# ---------------------------------------------------------------------------
full-stack: CURRENT_ENV=development
full-stack: check-docker check-env-files
	@echo "$(BLUE)Starting all services including monitoring...$(RESET)"
	docker compose --env-file .env.dev up --build -d
	@sleep 5
	@$(MAKE) health
	@echo "$(GREEN)All services started$(RESET)"
	@echo "API: http://api.localhost"
	@echo "Frontend: http://frontend.localhost"
	@echo "Mailhog: http://mailhog.localhost"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://grafana.localhost (admin/admin)"

# ---------------------------------------------------------------------------
# Monitoring: Start Prometheus and Grafana services
#
# This target starts the monitoring stack, including:
#   • Prometheus (metrics collection)
#   • Grafana (metrics visualization)
#
# Usage: make monitoring
# ---------------------------------------------------------------------------
monitoring:
	@echo "$(BLUE)Starting monitoring services (Prometheus and Grafana)...$(RESET)"
	docker compose --env-file .env up -d prometheus grafana
	@echo "$(GREEN)Monitoring services started$(RESET)"
	@echo "Prometheus UI: http://localhost:9090"
	@echo "Grafana UI: http://grafana.localhost (admin/admin)"

# ---------------------------------------------------------------------------
# Monitoring-Down: Stop Prometheus and Grafana services
#
# This target stops only the monitoring stack without affecting other services.
#
# Usage: make monitoring-down
# ---------------------------------------------------------------------------
monitoring-down:
	@echo "$(BLUE)Stopping monitoring services...$(RESET)"
	docker compose stop prometheus grafana
	docker compose rm -f prometheus grafana
	@echo "$(GREEN)Monitoring services stopped$(RESET)"

# -----------------------------------------
# Health Checks and Logs
# -----------------------------------------
health:
	docker compose ps

# Actively test the API health endpoint
health-api:
	@echo "Testing API health endpoint..."
	@curl -fsS http://localhost:8000/api/v1/health && echo "API healthy" || (echo "API unhealthy" && exit 1)

# Detailed health check including container status and API health endpoint test
health-detail:
	@echo "=== Container Status ==="
	@docker compose ps
	@echo "\n=== API Health Endpoint Check ==="
	@$(MAKE) health-api
	@echo "\n=== Docker Health Check Logs ==="
	@docker compose ps | grep -q "healthy" || (echo "⚠️ Unhealthy services detected" && exit 1)
	@echo "✅ All services healthy"

logs:
	docker compose logs

logs-frontend:
	docker compose logs frontend

logs-postgres:
	docker compose logs postgres

logs-redis:
	docker compose logs redis

logs-mailhog:
	docker compose logs mailhog

logs-prometheus:
	docker compose logs prometheus

logs-grafana:
	docker compose logs grafana

# -----------------------------------------
# Sentry Integration
# -----------------------------------------

# ---------------------------------------------------------------------------
# Verify-Sentry: Verify Backend Sentry configuration
#
# Checks if the SENTRY_BACKEND_DSN environment variable is set, which is required for Sentry to work.
#
# Usage: make verify-sentry
# ---------------------------------------------------------------------------
verify-sentry:
	@echo "$(BLUE)Verifying Sentry configuration...$(RESET)"
	@if [ -z "$$SENTRY_BACKEND_DSN" ]; then \
		echo "$(RED)Error: SENTRY_BACKEND_DSN environment variable is not set$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Sentry configuration verified.$(RESET)"

# ---------------------------------------------------------------------------
# Run-With-Sentry: Run with backend sentry enabled
#
# Runs API locally (not in Docker) with Sentry enabled, useful for testing Sentry integration during development.
#
# Usage: make run-with-sentry
# ---------------------------------------------------------------------------
run-with-sentry: verify-sentry
	@echo "$(BLUE)Running with Sentry enabled...$(RESET)"
	SENTRY_BACKEND_DSN=$(SENTRY_BACKEND_DSN) SENTRY_ENVIRONMENT=development go run ./cmd/api/main.go

# ---------------------------------------------------------------------------
# Test-Sentry: Test integration by triggering a test error
#
# Sends a test error to Sentry by making a request to a test endpoint, allowing verification that errors are being captured correctly.
#
# Usage: make test-sentry
# ---------------------------------------------------------------------------
test-sentry: verify-sentry
	@echo "$(BLUE)Sending test error to Sentry...$(RESET)"
	@curl -X GET http://api.localhost/api/v1/test-sentry || echo "\n$(GREEN)Test error sent to Sentry$(RESET)"

# ---------------------------------------------------------------------------
# Dev-With-Sentry: Start development environment with Sentry enabled
#
# Starts development environment with Sentry enabled, ensuring that the Sentry environment variables are set.
#
# Usage: make dev-with-sentry
# ---------------------------------------------------------------------------
dev-with-sentry: CURRENT_ENV=development
dev-with-sentry: check-docker check-env-files verify-sentry
	@echo "$(BLUE)Starting development environment with Sentry enabled...$(RESET)"
	docker compose --env-file .env.dev up --build -d
	@sleep 5
	@$(MAKE) health
	@echo "$(GREEN)Development environment with Sentry started$(RESET)"

# -----------------------------------------
# Troubleshooting
# -----------------------------------------
troubleshoot-frontend:
	@echo "=== Frontend Status ==="
	docker compose ps frontend
	@echo "\n=== Frontend Logs ==="
	docker compose logs --tail=50 frontend

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

troubleshoot-prometheus:
	@echo "=== Prometheus Status ==="
	docker compose ps prometheus
	@echo "\n=== Prometheus Logs ==="
	docker compose logs --tail=50 prometheus
	@echo "\n=== Prometheus Targets ==="
	@curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | {name: .labels.job, health: .health, lastError: .lastError}'

troubleshoot-grafana:
	@echo "=== Grafana Status ==="
	docker compose ps grafana
	@echo "\n=== Grafana Logs ==="
	docker compose logs --tail=50 grafana
	@echo "\n=== Grafana Health ==="
	@curl -s http://localhost:3000/api/health

# -----------------------------------------
# Help
# -----------------------------------------
help:
	@echo "$(BLUE)Available commands:$(RESET)"
	@echo "$(BLUE)Environment:$(RESET)"
	@echo " make init-env           - Initialize environment files from templates"
	@echo " make dev                - Start development environment (containerized)"
	@echo " make test               - Start test environment (containerized)"
	@echo " make prod               - Start production environment (containerized)"
	@echo " make restart-all        - Fully reset environment and start all services (backend & frontend)"
	@echo " make dev-backend        - Start development environment excluding the frontend (for backend development only)"
	@echo " make full-stack         - Start all services including monitoring stack"
	@echo " make down               - Shut down environment and remove volumes"
	@echo " make clean              - Remove all containers, volumes, and Docker artifacts"
	@echo "\n$(BLUE)Monitoring:$(RESET)"
	@echo " make monitoring         - Start Prometheus and Grafana services"
	@echo " make monitoring-down    - Stop Prometheus and Grafana services"
	@echo "\n$(BLUE)Health Checks:$(RESET)"
	@echo " make health             - Quick health status"
	@echo " make health-detail      - Detailed health status"
	@echo "\n$(BLUE)Logging:$(RESET)"
	@echo " make logs               - View all logs"
	@echo " make logs-postgres      - View Postgres logs"
	@echo " make logs-redis         - View Redis logs"
	@echo " make logs-mailhog       - View Mailhog logs"
	@echo " make logs-prometheus    - View Prometheus logs"
	@echo " make logs-grafana       - View Grafana logs"
	@echo "\n$(BLUE)Troubleshooting:$(RESET)"
	@echo " make troubleshoot-frontend - Troubleshoot Frontend"
	@echo " make troubleshoot-postgres - Troubleshoot Postgres"
	@echo " make troubleshoot-redis    - Troubleshoot Redis"
	@echo " make troubleshoot-prometheus  - Troubleshoot Prometheus"
	@echo " make troubleshoot-grafana     - Troubleshoot Grafana"
	@echo "\n$(BLUE)Sentry:$(RESET)"
	@echo " make verify-sentry      - Verify Sentry configuration"
	@echo " make run-with-sentry    - Run the API locally with Sentry enabled"
	@echo " make test-sentry        - Test Sentry integration by triggering a test error"
	@echo " make dev-with-sentry    - Start development environment with Sentry enabled"
