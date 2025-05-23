## Overview

Q-KO (pronounced "queue koh") is a "collector's companion" media management and billing tool designed for gamers. It tracks and manages metadata and associated expenses for physical and digital media collections.

## Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
  - [Modular Monolith Design](#modular-monolith-design)
- [Core Components](#core-components)
  - [Backend Architecture](#backend-architecture-go)
- [Scalability & Extensibility](#scalability--extensibility)
  - [Architectural Design for Scale](#architectural-design-for-scale)
  - [System Diagram](#system-diagram)
  - [Performance Optimizations](#performance-optimizations)
- [Data Flow](#data-flow)
  - [Request Handling](#request-handling)
  - [Business Logic](#business-logic)
  - [Data Access](#data-access)
  - [Caching Strategy](#caching-strategy)
- [Security Measures](#security-measures)
  - [Authentication & Authorization](#authentication--authorization-auth0)
  - [Application Security](#application-security-internal)
  - [API Security](#api-security-external)
- [Testing Strategy](#testing-strategy)
  - [Backend Testing](#backend-testing)
  - [Frontend Testing](#frontend-testing)
- [Monitoring and Observability](#monitoring-and-observability)
  - [Performance Metrics](#performance-metrics)
  - [Health Monitoring](#health-monitoring)
  - [Structured Logging](#structured-logging)
- [Development Workflow](#development-workflow)
  - [Version Control](#version-control)
  - [CI/CD Pipeline](#cicd-pipeline)
  - [Documentation](#documentation)

## System Architecture

### Modular Monolith Design

1. **Development Efficiency**
   - Single source of truth for code and configuration
   - Simplified dependency management
   - Unified tooling and development workflow
   - Centralized API contract management

2. **Technical Advantages**
   - Shared type definitions
   - Consistent API contracts
   - Unified domain models

## Core Components

The modular structure above enables Q-KO's scalable architecture.

### Backend Architecture (Go)

```text
.
backend/
├── cmd/                    # Application entrypoints
│   └── api/                # API server binary
│       └── main.go         # Entry point
├── config/                 # Configuration files
├── db/                     # Database related
│   ├── migrations/         # Database migrations
│   └── seeds/              # Seed data
├── server/                 # HTTP server package
│   ├── server.go           # Server definition
│   ├── routes.go           # Route configuration
│   └── tests/              # HTTP endpoint tests
└── internal/               # Private application code
    ├── books/              # Books domain module
    ├── core/               # Core business logic
    │   ├── models/         # Domain models
    │   ├── services/       # Business services
    │   └── repositories/   # Data access
    └── shared/             # Shared internal utilities
        ├── cache/          # Caching logic
        ├── crypto/         # Cryptography utilities
        ├── database/       # Database utilities
        ├── http/           # HTTP utilities
        ├── jwt/            # JWT handling
        ├── logger/         # Logging
        ├── middleware/     # Application middleware (including auth)
        ├── redis/          # Redis client
        ├── types/          # Common types
        └── validator/      # Input validation
```

### Database Architecture (PostgreSQL)

Q-KO uses PostgreSQL as its primary data store, chosen for its:

1. **Data Integrity**
   - ACID compliance for transaction reliability
   - Strong data consistency guarantees
   - Rich constraint system for data validation

2. **Performance Features**
   - Advanced indexing capabilities
   - Efficient query planning
   - Connection pooling support
   - Concurrent access optimization

3. **Scalability Support**
   - Partitioning for large datasets
   - Replication capabilities
   - JSON support for flexible schema evolution
   - Full-text search functionality

4. **Development Experience**
   - Migrations-based schema management
   - Comprehensive SQL feature set
   - Rich ecosystem of tools and extensions
   - Strong type system alignment with Go

### Frontend Architecture (React)

```text
frontend/
├── src/
│   ├── core/                    # Core application infrastructure
│   │   ├── api/                 # API integration layer
│   │   │   ├── client/          # Axios client configuration
│   │   │   ├── services/        # API service definitions
│   │   │   ├── queries/         # TanStack Query hooks
│   │   │   ├── adapters/        # Data transformation layer
│   │   │   ├── hooks/           # Custom API hooks
│   │   │   └── types/           # API type definitions
│   │   ├── auth/                # Authentication management
│   │   ├── error/               # Error handling
│   │   ├── network-status/      # Network state management
│   │   └── theme/               # Theme configuration
│   ├── features/                # Feature-specific modules
│   │   ├── library/             # Library management
│   │   ├── wishlist/            # Wishlist functionality
│   │   ├── search/              # Search interface
│   │   └── analytics/           # Analytics dashboard
│   ├── shared/                  # Shared components and utilities
│   │   ├── components/          # Reusable UI components
│   │   ├── hooks/               # Shared React hooks
│   │   └── utils/               # Utility functions
│   └── types/                   # Global type definitions
└── public/                      # Static assets
```

The frontend architecture emphasizes maintainability, performance, and developer experience through carefully selected technologies and patterns:

1. **API Integration Layer**
   - Axios client for HTTP requests with interceptors and error handling
   - TanStack Query for server state management and caching
   - Type-safe API services with TypeScript
   - Data adapters for transforming API responses to UI models
   - Custom hooks for reusable API logic

2. **State Management**
   - TanStack Query for server state
   - React Context for global UI state
   - Local component state for UI-specific concerns
   - Optimistic updates for improved UX
   - Automatic background refetching and cache invalidation

3. **UI Component Strategy**
   - ShadCN UI + Tailwind for consistent, modern UI
   - Component composition for reusability
   - Responsive design patterns
   - Accessibility-first approach
   - Performance-optimized rendering

4. **Feature Organization**
   - Domain-driven feature modules
   - Shared core infrastructure
   - Clear separation of concerns
   - Scalable module structure
   - Type-safe feature boundaries

5. **Development Experience**
   - TypeScript for type safety
   - ESLint for code quality
   - Vite for fast development
   - Hot module replacement
   - Development proxy configuration

6. **Build and Deployment**
   - Vite for optimized builds
   - Environment-specific configurations
   - Docker containerization
   - Static asset optimization
   - Development and production builds

This architecture supports Q-KO's requirements for a modern, maintainable, and performant frontend application while providing a great developer experience.

## Scalability & Extensibility

These scalability patterns inform our data flow design:

### Architectural Design for Scale

1. **Modular Domain Extension**
   - Designed for multi-domain support (books, music, movies)
   - Domain-specific modules with shared core infrastructure
   - Factory pattern initialization enables seamless domain addition
   ```go
   // Example domain registration in factory.go
   bookDomainHandler := domains.NewBookDomainHandler(bookHandlers, log)
   operationsManager.RegisterDomain(bookDomainHandler)
   ```

2. **Layered Service Architecture**
   ```
   Request Flow:
   Handler → Service → Operations → Repository
   └── Validation
   └── Caching
   └── Authentication
   ```

3. **Resource Management**
   - Graceful shutdown coordination
   - Connection pooling for database access
   - Redis caching with invalidation strategies
   - Worker pool management for background tasks

### System Diagram
```mermaid
graph TD
    Client[Client] --> Router[Chi Router/Middleware]

    subgraph API Layer
        Router --> Auth[Auth Handler]
        Router --> Library[Library Handlers]
        Router --> Wishlist[Wishlist Handlers]
        Router --> Search[Search Handlers]
        Router --> Analytics[Analytics Handlers]
        Router --> Health[Health Handlers]
    end

    subgraph Services
        Library --> LibraryService[Library Service]
        Wishlist --> WishlistService[Wishlist Service]
        Search --> SearchService[Search Service]
        Analytics --> AnalyticsService[Analytics Service]
    end

    subgraph Infrastructure
        LibraryService --> IGDB[IGDB Integration]
        LibraryService --> MediaStorage[Media Storage]
        LibraryService --> LocationService[Location Service]
    end

    subgraph Operations
        LibraryService --> OpManager[Operations Manager]
        WishlistService --> OpManager
        SearchService --> OpManager
        AnalyticsService --> OpManager
        OpManager --> Cache[(Redis Cache)]
        OpManager --> DB[(PostgreSQL)]
    end

    subgraph Monitoring
        Health --> HealthService[Health Service]
        HealthService --> Monitoring[Monitoring System]
        Monitoring --> Metrics[Performance Metrics]
        Monitoring --> Logging[Structured Logging]
    end

    subgraph Background
        Workers[Background Workers]
        Cache --> Workers
        DB --> Workers
        AnalyticsService --> Workers
    end
```

The system diagram above illustrates Q-KO's comprehensive architecture, organized into several key layers:

1. **API Layer**
   - Handles all incoming HTTP requests through Chi Router
   - Implements middleware for authentication, logging, and request validation
   - Routes requests to appropriate domain handlers
   - Includes dedicated handlers for library, wishlist, search, analytics, and health monitoring

2. **Services Layer**
   - Implements core business logic for each domain
   - Maintains separation of concerns between different features
   - Provides clean interfaces for domain operations
   - Coordinates between infrastructure and domain logic

3. **Infrastructure Layer**
   - Integrates with external services (IGDB for game metadata)
   - Manages media storage for collection items
   - Handles location tracking and management
   - Provides shared infrastructure services

4. **Operations Layer**
   - Manages data persistence through PostgreSQL
   - Implements caching strategies with Redis
   - Coordinates database operations and transactions
   - Handles data consistency and integrity

5. **Monitoring Layer**
   - Provides health check endpoints
   - Collects and processes performance metrics
   - Implements structured logging (Not yet implemented)
   - Enables system observability  (Not yet implemented)

6. **Background Processing**
   - Handles asynchronous tasks
   - Processes analytics data
   - Manages cache invalidation
   - Coordinates background jobs

This architecture supports Q-KO's requirements for scalability, maintainability, and extensibility while providing a clear separation of concerns between different system components.

### Performance Optimizations

1. **Multi-Level Caching**
   - Redis for distributed caching
   - In-memory caching for frequent operations
   - Cache invalidation workers for consistency

2. **Request Processing**
   - Rate limiting by operation type
   ```go
   r.With(middleware.StandardRateLimiter).Get("/books")
   r.With(middleware.IntensiveRateLimiter).Get("/search")
   ```
   - Adaptive compression for large responses
   - Request validation pipeline

3. **Resource Management**
   - Connection pooling
   - Prepared statement caching
   - Graceful shutdown coordination
   ```go
   func gracefulShutdown(ctx context.Context, srv *http.Server, f *Factory, app *application, log *slog.Logger) error {
       // Coordinated shutdown sequence
       // 1. Stop accepting new requests
       // 2. Complete in-flight requests
       // 3. Close background workers
       // 4. Close database connections
   }
   ```

## Data Flow

1. **Request Handling**
   - JWT-based authentication (Not yet implemented, Auth0 integration planned)
   - Request validation
   - Rate limiting
   - CORS policy enforcement

2. **Business Logic**
   - Domain-driven design principles
   - Service layer abstraction
   - Strong type safety

3. **Data Access**
   - Repository pattern
   - Connection pooling
   - Prepared statements
   - Transaction management

4. **Caching Strategy**
   - Redis for frequently accessed data
   - Cache invalidation policies
   - Cache-aside pattern implementation

## Security Measures

### Authentication & Authorization (Auth0)

1. **Identity & Access Management** (Not yet implemented)
   - Frontend authentication via Auth0 SPA SDK
   - Backend validation via Auth0 JWT middleware
   - OAuth 2.0 and OpenID Connect compliance
   - Social login integration

2. **Security Features**
   - MFA support
   - Secure session management with refresh tokens (frontend)
   - Role-based access control (backend)
   - Token validation and scope checking (backend)

### Application Security (Internal)

1. **Request Protection**
   - CSRF protection via gorilla/csrf
   - Rate limiting
   - Input validation and sanitization
   - Secure headers management via chi middleware

2. **Data Security**
   - SQL injection prevention through prepared statements (Not yet implemented)
   - Audit logging of authentication events               (Not yet implemented)
   - Sensitive data encryption                            (Not yet implemented)
   - Error handling security

### API Security (External)

1. **Access Control**
   - Secure cookie handling
   - API rate limiting
   - Request validation middleware
   - CORS policy enforcement via chi middleware

2. **API Hardening**
   - Request throttling
   - API versioning
   - Response sanitization
   - Security headers enforcement

## Testing Strategy

### Backend Testing

1. **Test Types**
   - Unit tests with `testify`
   - Integration tests
   - API endpoint tests
   - Performance benchmarks

### Frontend Testing (Currently in Planning)

1. **Test Types**
   - Unit tests with Jest
   - Component tests with React Testing Library


## Monitoring and Observability

1. **Performance Metrics**
   - Cache hit/miss ratios
   - Request latencies
   - Worker queue depths
   ```go
   metrics := f.CacheManager.GetMetrics()
   log.Info("Cache metrics",
       "totalOps", metrics.TotalOps,
       "l1Failures", metrics.L1Failures,
       "l2Failures", metrics.L2Failures,
   )
   ```

2. **Health Monitoring**
   - Component health checks
   - Resource utilization tracking
   - Error rate monitoring

3. **Structured Logging**
   - Request tracing
   - Error tracking
   - Performance monitoring

## Development Workflow

1. **Version Control**
   - Feature branch workflow
   - Conventional commits
   - Pull request templates
   - Code review guidelines

2. **CI/CD Pipeline**
   - Automated testing
   - Code quality checks
   - Security scanning
   - Automated deployments

3. **Documentation**
   - API documentation (OpenAPI/Swagger)
   - Component documentation
   - Architecture decisions records (ADRs)