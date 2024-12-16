#=======================================================
# QKO BETA API - Development Operations
# Version: 1.0.0
#
# Required:
# - Docker 20.10.0+
# - Docker Compose 2.0.0+
#=======================================================

# Declare all phony targets
.PHONY: dev test down health health-detail logs logs-postgres logs-redis logs-mailhog troubleshoot-postgres troubleshoot-redis help clean

# Colors for output
BLUE := \033[34m
GREEN := \033[32m
RED := \033[31m
RESET := \033[0m

# Check required tools
check-docker:
	@if [ -z "$(shell docker-compose --version 2>/dev/null)" ]; then \
		echo "$(RED)Error: docker-compose is not installed$(RESET)"; \
		exit 1; \
	fi

# Check environment files
check-env-files:
	@if [ ! -f .env.dev ]; then \
		echo "$(RED)Error: .env.dev not found. Copy .env.example to .env.dev and configure$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f .env.test ]; then \
		echo "$(RED)Error: .env.test not found. Copy .env.example to .env.test and configure$(RESET)"; \
		exit 1; \
	fi

#==========================================
# Environment Management
#==========================================

# Start development environment
dev: check-docker check-env-files
	docker-compose --env-file .env.dev up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@make health

# Start test environment
test: check-docker check-env-files
	docker-compose --env-file .env.test up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@make health

# Shut down environment and remove volumes
down:
	docker-compose down -v

#==========================================
# Cleanup
#==========================================

# Remove all containers, volumes, and docker artifacts
clean:
	@echo "Cleaning up all resources..."
	docker-compose down -v
	docker system prune -f
	@echo "$(GREEN)Cleanup complete$(RESET)"

#==========================================
# Health Checks
#==========================================

# Quick health status check
health:
	docker-compose ps

# Detailed health status with logs
health-detail:
	@echo "=== Container Status ==="
	docker-compose ps
	@echo "\n=== Health Check Logs ==="
	docker-compose ps | grep -q "healthy" || (echo "⚠️  Unhealthy services detected" && exit 1)
	@echo "✅ All services healthy"

#==========================================
# Logging
#==========================================

# View all service logs
logs:
	docker-compose logs

# View Postgres logs
logs-postgres:
	docker-compose logs postgres

# View Redis logs
logs-redis:
	docker-compose logs redis

# View Mailhog logs
logs-mailhog:
	docker-compose logs mailhog

#==========================================
# Troubleshooting
#==========================================

# Detailed Postgres troubleshooting
troubleshoot-postgres:
	@echo "=== Postgres Status ==="
	docker-compose ps postgres
	@echo "\n=== Postgres Logs ==="
	docker-compose logs --tail=50 postgres
	@echo "\n=== Postgres Environment ==="
	docker-compose exec postgres env

# Detailed Redis troubleshooting
troubleshoot-redis:
	@echo "=== Redis Status ==="
	docker-compose ps redis
	@echo "\n=== Redis Logs ==="
	docker-compose logs --tail=50 redis
	@echo "\n=== Redis Ping Test ==="
	docker-compose exec redis redis-cli ping

#==========================================
# Help
#==========================================

# Show this help
help:
	@echo "$(BLUE)Available commands:$(RESET)"
	@echo "$(BLUE)Environment:$(RESET)"
	@echo "  make dev              - Start development environment"
	@echo "  make test             - Start test environment"
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