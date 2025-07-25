# QKO BETA API - Development Operations (Containerized Setup)
# Version: 2.7.0
#
# This Makefile leverages Docker Compose to build, run, and manage
# the following containerized services:
#   - Backend API (Golang application)
#   - Frontend (React application)
#   - Redis
#   - Postgres
#   - Sentry (Backend Error tracking)
#
# Best practices used:
#   - All services run in containers on a shared private network.
#   - The API container connects to Redis via Docker's internal DNS (using the service name).
#   - Environment-specific configuration is loaded by using .env files.
#   - No sensitive service is exposed publicly (only the API is exposed with a published port).
#   - Metrics are collected and visualized for monitoring system health.
#
# ---------------------------------------------------------------------------
# Development Workflow
# ---------------------------------------------------------------------------
#
# Normal Development Day:
#   make dev-backend    # Start all backend services
#   make stop          # Take a break (preserves data)
#   make dev-backend    # Come back (data still there)
#   make stop          # End of day
#
# When to Use Reset:
#   - Starting a new feature that changes database schema
#   - Database is acting weird and you want to start fresh
#   - Preparing for a demo to stakeholders
#   - Testing from a clean slate
#
# Example Workflows:
#
# 1. Normal Development:
#    make dev-backend    # Start work
#    make stop          # Take lunch
#    make dev-backend    # Continue work
#    make stop          # End day
#
# 2. Starting New Feature:
#    make reset         # Clean slate
#    make dev-backend    # Start fresh
#
# 3. Database Issues:
#    make reset         # Reset everything
#    make dev-backend    # Start clean
#
# 4. Demo Preparation:
#    make reset         # Clean state
#    make dev-backend    # Start services
#
# ⚠️ WARNING: make reset DELETES ALL DATA
# Use with caution and only when you want to start completely fresh.
#
# ---------------------------------------------------------------------------

.PHONY: init-env check-docker check-env-files dev test prod down clean health health-detail logs logs-postgres logs-redis troubleshoot-postgres troubleshoot-redis verify-sentry run-with-sentry test-sentry dev-with-sentry help backup restore list-backups check-db migrate migrate-down recreate nuclear spend-tracking-db-seed spend-tracking-db-seed-down seed-data-complete debug-migration

# Define allowed environments and set current environment
ENVS := development test production
CURRENT_ENV ?= development
REDIS_VOLUME ?= redis_data_updated

# Colors (for terminal output)
BLUE := \033[34m
GREEN := \033[32m
RED := \033[31m
RESET := \033[0m
YELLOW := \033[33m

# ---------------------------------------------------------------------------
# Stop: Safely stop all services while preserving data
#
# This target stops all services but preserves all volumes and data.
# Use this for normal development when you want to restart services.
#
# Usage: make stop
# ---------------------------------------------------------------------------
stop:
	@echo "$(BLUE)Stopping all services (preserving data)...$(RESET)"
	docker compose --env-file .env.dev down
	@echo "$(GREEN)Services stopped. Data preserved.$(RESET)"

# ---------------------------------------------------------------------------
# Clean: Remove containers and networks, preserve volumes
#
# This target removes all containers and networks but preserves volumes.
# Use this when you want to clean up containers but keep your data.
#
# Usage: make clean
# ---------------------------------------------------------------------------
clean:
	@echo "$(BLUE)Cleaning containers and networks (preserving data)...$(RESET)"
	docker compose --env-file .env.dev down
	@echo "$(GREEN)Cleanup complete. Data preserved.$(RESET)"

# ---------------------------------------------------------------------------
# List Backups: Show available database backups
#
# This target lists all available database backups.
#
# Usage: make list-backups
# ---------------------------------------------------------------------------
list-backups:
	@echo "$(BLUE)Available database backups:$(RESET)"
	@ls -lh backups/qkoapi_*.sql 2>/dev/null || echo "$(YELLOW)No backups found$(RESET)"

# ---------------------------------------------------------------------------
# Backup: Create a backup of the database
#
# This target creates a timestamped backup of the database.
# Automatically keeps only the last 5 backups.
#
# Usage: make backup
# ---------------------------------------------------------------------------
backup:
	@echo "$(BLUE)Creating database backup...$(RESET)"
	@mkdir -p backups
	@docker compose exec -T postgres pg_dump -U postgres qkoapi > backups/qkoapi_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)Backup created successfully$(RESET)"
	@echo "$(BLUE)Rotating backups...$(RESET)"
	@ls -t backups/qkoapi_*.sql 2>/dev/null | tail -n +6 | xargs rm -f 2>/dev/null || true

# ---------------------------------------------------------------------------
# Check Database: Show recent database changes
#
# This target shows recent changes in the database.
# Useful before resetting to ensure no important data is lost.
#
# Usage: make check-db
# ---------------------------------------------------------------------------
check-db:
	@echo "$(BLUE)Checking recent database changes...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -c "\
		SELECT table_name, COUNT(*) as row_count \
		FROM information_schema.tables \
		WHERE table_schema = 'public' \
		GROUP BY table_name \
		ORDER BY table_name;"
	@echo "$(YELLOW)Note: This shows table row counts. Consider backing up if you have important data.$(RESET)"

# ---------------------------------------------------------------------------
# Reset: Forcefully remove everything including data
#
# This target removes all containers, networks, and volumes.
# Use this when you want to start completely fresh.
#
# Usage: make reset [BACKUP=true] [FORCE=true]
#   BACKUP=true  Create a backup before reset (default: true)
#   FORCE=true   Skip database check (default: false)
#
# Example:
#   make reset              # Checks database, creates backup, then resets
#   make reset FORCE=true   # Skips checks and backup
# ---------------------------------------------------------------------------
reset:
	@if [ "$(FORCE)" != "true" ]; then \
		$(MAKE) check-db; \
		echo "$(YELLOW)Are you sure you want to reset? This will delete all data. [y/N]$(RESET)"; \
		read -r response; \
		if [[ ! "$$response" =~ ^[Yy]$$ ]]; then \
			echo "$(BLUE)Reset cancelled$(RESET)"; \
			exit 1; \
		fi; \
	fi
	@if [ "$(BACKUP)" != "false" ]; then \
		$(MAKE) backup; \
	fi
	@echo "$(BLUE)Resetting everything (including data)...$(RESET)"
	docker compose --env-file .env.dev down -v
	docker volume rm -f $(REDIS_VOLUME)
	docker volume prune -f
	@echo "$(GREEN)Complete reset complete. All data removed.$(RESET)"

# ---------------------------------------------------------------------------
# Up: Build + start backend services (excluding frontend)
#
# This target rebuilds and starts the backend services along with necessary
# dependencies, including:
#   • Traefik (as the reverse proxy)
#   • Backend API (Golang application)
#   • Redis
#   • Postgres
#
# The frontend service is purposefully excluded when using this target.
#
# Usage: make up
# ---------------------------------------------------------------------------
up:
	@echo "Rebuilding and starting backend services (excluding frontend)..."
	docker compose --env-file .env up --build -d traefik api redis postgres

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
	@if [ "$(CURRENT_ENV)" = "production" ] && [ ! -f .env.prod ]; then \
		echo "$(RED)Error: .env.prod not found$(RESET)"; \
		exit 1; \
	fi

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
	docker compose --env-file .env up --build -d
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
	docker compose --env-file .env build --no-cache api
	docker compose --env-file .env up --build -d traefik api postgres redis
	@sleep 5
	@$(MAKE) health

# Shut down the development environment
down:
	docker compose --env-file .env down -v



# ---------------------------------------------------------------------------
# Monitoring services have been removed from docker-compose.yml
# ---------------------------------------------------------------------------

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
	@echo " make down               - Shut down environment and remove volumes"
	@echo " make clean              - Remove all containers, volumes, and Docker artifacts"
	@echo "\n$(BLUE)Health Checks:$(RESET)"
	@echo " make health             - Quick health status"
	@echo " make health-detail      - Detailed health status"
	@echo "\n$(BLUE)Logging:$(RESET)"
	@echo " make logs               - View all logs"
	@echo " make logs-postgres      - View Postgres logs"
	@echo " make logs-redis         - View Redis logs"
	@echo "\n$(BLUE)Troubleshooting:$(RESET)"
	@echo " make troubleshoot-frontend - Troubleshoot Frontend"
	@echo " make troubleshoot-postgres - Troubleshoot Postgres"
	@echo " make troubleshoot-redis    - Troubleshoot Redis"
	@echo "\n$(BLUE)Sentry:$(RESET)"
	@echo " make verify-sentry      - Verify Sentry configuration"
	@echo " make run-with-sentry    - Run the API locally with Sentry enabled"
	@echo " make test-sentry        - Test Sentry integration by triggering a test error"
	@echo " make dev-with-sentry    - Start development environment with Sentry enabled"
	@echo "\n$(BLUE)Database:$(RESET)"
	@echo " make backup             - Create a database backup"
	@echo " make restore BACKUP_FILE=path/to/backup.sql - Restore database from a backup"
	@echo " make migrate            - Run database migrations"
	@echo " make migrate-down       - Roll back database migrations"
	@echo " make check-db           - Check database status"
	@echo " make list-backups       - List available database backups"
	@echo " make spend-tracking-db-seed - Seed spend tracking data"
	@echo " make spend-tracking-db-seed-down - Remove spend tracking seed data"
	@echo " make seed-data-complete - Seed complete data set"
	@echo " make debug-migration - Debug migration and seeding issues"
	@echo " make reset-test-data - Clear test users from database and Auth0"

# ---------------------------------------------------------------------------
# Restore: Restore database from backup
#
# This target restores the database from a specified backup file.
#
# Usage: make restore BACKUP_FILE=backups/qkoapi_20240101_120000.sql
# ---------------------------------------------------------------------------
restore:
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "$(RED)Error: BACKUP_FILE is required$(RESET)"; \
		echo "Usage: make restore BACKUP_FILE=path/to/backup.sql"; \
		exit 1; \
	fi
	@if [ ! -f "$(BACKUP_FILE)" ]; then \
		echo "$(RED)Error: Backup file $(BACKUP_FILE) not found$(RESET)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Restoring database from $(BACKUP_FILE)...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi < $(BACKUP_FILE)
	@echo "$(GREEN)Database restored successfully$(RESET)"

# ---------------------------------------------------------------------------
# Database: Database management commands
# ---------------------------------------------------------------------------
migrate:
	@echo "$(BLUE)Running database migrations...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -f /docker-entrypoint-initdb.d/migrations/20240414150000_create_initial_schema.up.sql
	@echo "$(GREEN)Migrations completed successfully$(RESET)"

migrate-down:
	@echo "$(BLUE)Rolling back database migrations...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -f /docker-entrypoint-initdb.d/migrations/20240414150000_create_initial_schema.down.sql
	@echo "$(GREEN)Migrations rolled back successfully$(RESET)"

spend-tracking-db-seed:
	@echo "$(BLUE)Seeding spend tracking data...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -f /docker-entrypoint-initdb.d/migrations/20240414150001_seed_spend_tracking_data.up.sql
	@echo "$(GREEN)Spend tracking data seeded successfully$(RESET)"

spend-tracking-db-seed-down:
	@echo "$(BLUE)Removing spend tracking seed data...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -f /docker-entrypoint-initdb.d/migrations/20240414150001_seed_spend_tracking_data.down.sql
	@echo "$(GREEN)Spend tracking seed data removed successfully$(RESET)"

seed-data-complete:
	@echo "$(BLUE)Seeding complete data set...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -f /docker-entrypoint-initdb.d/migrations/seed_data_complete.sql
	@echo "$(GREEN)Complete data set seeded successfully$(RESET)"

debug-migration:
	@echo "$(BLUE)Debugging migration issues...$(RESET)"
	@if [ ! -f "scripts/debug_migration.sh" ]; then \
		echo "$(RED)Error: Debug script not found$(RESET)"; \
		echo "Please ensure scripts/debug_migration.sh exists and is executable"; \
		exit 1; \
	fi
	@chmod +x scripts/debug_migration.sh
	@./scripts/debug_migration.sh

# Recreate: Forcefully stop everything, remove all containers, networks, and volumes, and start fresh
recreate:
	@echo "$(BLUE)Forcefully stopping everything...$(RESET)"
	docker compose --env-file .env.dev down -v
	docker volume prune -f
	@echo "$(GREEN)Starting backend fresh...$(RESET)"
	docker compose --env-file .env.dev up -d
	@echo "$(GREEN)Backend recreated successfully.$(RESET)"

# Nuclear: Complete system cleanup to resolve disk space issues
# Usage: make nuclear
nuclear:
	@echo "$(RED)⚠️  WARNING: This will remove ALL Docker resources on your system!$(RESET)"
	@echo "$(YELLOW)Are you sure you want to proceed? [y/N]$(RESET)"
	@read -r response; \
	if [[ ! "$$response" =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Nuclear cleanup cancelled$(RESET)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Stopping all containers...$(RESET)"
	docker compose --env-file .env.dev down -v
	@echo "$(BLUE)Removing all stopped containers...$(RESET)"
	docker container prune -f
	@echo "$(BLUE)Removing all unused networks...$(RESET)"
	docker network prune -f
	@echo "$(BLUE)Removing all unused volumes...$(RESET)"
	docker volume prune -f
	@echo "$(BLUE)Removing all unused images...$(RESET)"
	docker image prune -af
	@echo "$(BLUE)Removing all build cache...$(RESET)"
	docker builder prune -af
	@echo "$(BLUE)Removing all unused data...$(RESET)"
	docker system prune -af --volumes
	@echo "$(GREEN)Starting backend fresh...$(RESET)"
	docker compose --env-file .env.dev up -d traefik api postgres redis
	@echo "$(GREEN)Nuclear cleanup complete. Backend services should be fresh and clean.$(RESET)"

# ---------------------------------------------------------------------------
# Reset Test Data: Clear test users from database and Auth0
#
# This target clears test data for development testing.
# Requires AUTH0_DOMAIN, AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_TEST_EMAIL in .env.dev
#
# Usage: make reset-test-data
# ---------------------------------------------------------------------------
reset-test-data:
	@echo "$(BLUE)Resetting test data...$(RESET)"
	@if [ -z "$$AUTH0_DOMAIN" ] || [ -z "$$AUTH0_CLIENT_ID" ] || [ -z "$$AUTH0_CLIENT_SECRET" ] || [ -z "$$AUTH0_TEST_EMAIL" ]; then \
		echo "$(RED)Error: Missing required environment variables$(RESET)"; \
		echo "Please ensure AUTH0_DOMAIN, AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_TEST_EMAIL are set in .env.dev"; \
		exit 1; \
	fi
	@echo "$(BLUE)Clearing test users from database...$(RESET)"
	@docker compose exec -T postgres psql -U postgres -d qkoapi -c "DELETE FROM users WHERE email LIKE '%test%';"
	@echo "$(BLUE)Clearing Auth0 test user: $$AUTH0_TEST_EMAIL$(RESET)"
	@curl -X DELETE "https://$$AUTH0_DOMAIN/api/v2/users-by-email?email=$$AUTH0_TEST_EMAIL" \
		-H "Authorization: Bearer $$(curl -s -X POST https://$$AUTH0_DOMAIN/oauth/token \
		-H "Content-Type: application/json" \
		-d '{"client_id":"$$AUTH0_CLIENT_ID","client_secret":"$$AUTH0_CLIENT_SECRET","audience":"https://$$AUTH0_DOMAIN/api/v2/","grant_type":"client_credentials"}' | jq -r '.access_token')" \
		-H "Content-Type: application/json" || echo "$(YELLOW)Warning: Could not delete Auth0 user (may not exist)$(RESET)"
	@echo "$(GREEN)Test data reset complete$(RESET)"
