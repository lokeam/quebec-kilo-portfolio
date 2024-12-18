# API Documentation

## Overview
Q-KO's API is designed following OpenAPI 3.0 specification standards, demonstrating API-first development practices while maintaining simplicity for the current project scale.

## Development Approach
While enterprise projects often use automated code generation from OpenAPI specs (especially useful for large-scale APIs with multiple services), this project implements endpoints manually for:

- Better learning experience
- More direct control over implementation
- Appropriate complexity for current scale

## Structure
```text
docs/api/
├── _working/
│   └── current.yaml
├── openapi.yaml
├── paths/
│   └── health.yaml
└── components/
    ├── responses/          # Reusable responses
    │   ├── error.yaml
    │   └── not-found.yaml
    ├── schemas/           # Data models
    │   ├── user.yaml
    │   └── book.yaml
    ├── parameters/        # URL parameters
    │   ├── user-id.yaml
    │   └── book-id.yaml
    └── securitySchemes/   # Auth methods
        ├── bearer-auth.yaml
        └── api-key.yaml
```

## Development Workflow
1. Design new endpoints in `_working/current.yaml` using Swagger Editor
2. Once validated, split into appropriate files:
   - Endpoint definitions → `paths/`
   - Common responses → `components/responses/`
   - Security schemes → `components/security/`
   - Data models → `schemas/`
3. Update `openapi.yaml` with new references
4. Commit changes with descriptive messages

## Security
- Authentication using JWT
- CSRF protection
- Rate limiting by operation type
- Input validation pipeline

## Endpoints
Documentation for all API endpoints can be found in the OpenAPI specification files.