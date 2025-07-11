# Q-KO Beta

🚧 **Project Status:** Late Stage Implementation

**Q-KO** is a comprehensive inventory and media management tool. It helps users track their physical and digital media assets, associate metadata like creators, genres, and tags, and access those items through various platforms and online services.

This mono repo project demonstrating enterprise-level software practices through:

- Clean Architecture principles
- Domain-Driven Design
- Production-grade security
- Scalable infrastructure
- Comprehensive testing

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Development](#development)
- [Testing](#testing)
- [CI/CD](#cicd)

## Features

### Core Features
- **Unified Item Management:** Store and categorize your library of either video games or books with rich metadata curated from the Video Game DB and Google BooksAPIs.
- **Creators & Genres:** Associate creators and genres to items for better discovery.
- **Online Service Expense Tracking:** Track monthly and yearly expenses associated with online services.
- **Physical & Digital Locations:** Track where items are stored physically and where they can be accessed online.
- **User Ownership & Wishlists:** Allow users to maintain collections and wishlists, making the tool ideal for collectors or agencies managing multiple media assets.

### Technical Features
- **Modular Monolith:** Built for easy feature addition and future service separation
- **API-First Design:** REST API documented with OpenAPI/Swagger
- **Caching Strategy:** Optimized performance with Redis
- **Secure by Default:** JWT authentication, RBAC, and input validation via Auth0, Go-Chi, Go-Validator, and Go-JWT

## Tech Stack

### Infrastructure
- **Docker** (20.10.0+) - Containerization and development environment
- **Make** - Build automation and development workflow

### Backend
- **Golang** (1.21) - Server-side application
- **PostgreSQL** (14) - Primary data store
- **Redis** (7.0) - Caching layer

### Frontend
- **TypeScript** (5.2.2) - Type-safe JavaScript
- **React** (18.3.1) - UI framework
- **Vite** (5.3.1) - Build tool and dev server

## API Documentation
Q-KO's API is designed and documented using OpenAPI 3.0 specification:
- [Interactive API Documentation](./docs/api) - OpenAPI/Swagger UI
- [Architecture Documentation](./docs/architecture.md) - Detailed system design
- [API Reference](./docs/api/openapi.yaml) - OpenAPI specification

## Architecture
See [Architecture Documentation](docs/architecture.md)

## Setup & Installation

### Prerequisites
- Docker (20.10.0+)
- Docker Compose (2.0.0+)
- Make

### Environment Setup

## Database
This project uses PostgreSQL for data persistence. The database is configured for security and performance in both development and production environments.

For detailed information about:
- Database configuration
- Development environment setup
- Access methods
- Migration procedures

See [database.md](docs/technical/database.md) (for team members only - not included in public repositories).

Quick development access:
- Command line: `docker compose exec postgres psql -U postgres -d qkoapi`

## Redis Setup & Configuration
See [Redis Configuration](docs/technical/redis.md) for more information.

## Testing

Testing strategy and instructions coming soon. Will include:
- Unit tests
- Integration tests
- End-to-end tests
- Performance testing

## CI/CD

CI/CD pipeline documentation coming soon. Will include:
- Automated testing
- Code quality checks
- Build process
- Deployment strategy
