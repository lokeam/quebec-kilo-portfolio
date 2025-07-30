
<img width="1268" height="902" alt="mkt_page" src="https://github.com/user-attachments/assets/77b819d9-f1b3-4204-8145-e17de4dcce01" />

# QKO Game Management App

[![Beta Access](https://img.shields.io/badge/Beta-Access-orange)](https://www.q-ko.com)\
<a href="https:www.q-ko.com">https://www.q-ko.com</a>

**QKO** is an end‑to‑end inventory and media management platform:


* Track everything about your physical & digital video game collection
* Attach rich metadata (creators, genres, purchase history)
* Automate expense reporting for subscriptions
* Access via web, mobile‑optimized UI, and REST API

## Screenshots

<img width="334" alt="library_detail_card" src="https://github.com/user-attachments/assets/05c22e7c-ee49-437f-ae7d-c80a822ba738" />
<img width="334" alt="add_game_to_library_form_2" src="https://github.com/user-attachments/assets/254666a0-88fc-4a10-9a3a-f7142566e425" />

<img width="668" alt="physical_locations" src="https://github.com/user-attachments/assets/7edd765a-95a6-4805-94b5-e6f3f470e99b" />
<img width="668" alt="spend_tracking" src="https://github.com/user-attachments/assets/1bc26b47-8950-45f0-8a79-57fb56960e50" />

## Table of Contents

1. [Features](#features)
2. [Technical Highlights](#technical-highlights)
3. [Documentation & References](./docs)
4. [Getting Started (Beta)](#getting-started-beta)
5. [Future Roadmap](#roadmap)
6. [License](#license)

## Features

### User‑Facing Features
- **Unified Collections**: single view for physical and digital assets
- **Smart Tags & Metadata**: auto‑fetch from VideoGameDB
- **Expense Tracking**: monthly/yearly spend analytics for subscriptions
- **Multi‑Platform Access**: web UI, mobile web, and REST API


## Technical Highlights
<details>
  <summary>🛠️ Tech Stack & Tooling</summary>

  **Infrastructure**
  - Docker (20.10+), Docker Compose (2.0+), Make
  - GitHub Actions → Docker → Cloudflare Pages/Workers

  **Backend**
  - Go (1.21), Chi, go-validator
  - PostgreSQL 14, Redis 7.0

  **Frontend**
  - React 18 + Vite 5
  - TypeScript 5

  **Marketing Site**
  - Astro 5 + Tailwind
</details>

<details>
  <summary>📦 Modular Monolith with Clean Architecture</summary>

  **What**:
  * The entire codebase lives in a single repository but is strictly layered:

  1. Domain (core business models & interfaces)
  2. Application (use‑case orchestration)
  3. Adapters (HTTP, DB, cache, external APIs)

  **How**:
  * Each layer depends only on the layer immediately beneath it. We use Go’s module system to enforce clear import boundaries, and interfaces to decouple business logic from transport or persistence.

  **Why**:
  * You get the simplicity of a monolith for local development and CI‑runs, plus the clean seams you need to slice off services into microservices or serverless functions later.
</details>


<details>
  <summary>🔗 API‑First Design</summary>

  **What**:
  * Every endpoint is defined in an OpenAPI 3.0 spec.

  **How**:
  * We use oapi-codegen to generate both server stubs and TypeScript/Go clients directly from openapi.yaml. Our CI job runs a diff on the spec to catch accidental breaking changes.

  **Why**:
  * Guarantees end‑to‑end contract consistency, speeds up frontend integration, and provides interactive Swagger UI for internal and beta consumers.
</details>

<details>
  <summary>⚡ Performance Optimizations</summary>

  **Redis Caching**:
  * Query‑level caching with 5‑minute TTL for “read‐heavy” endpoints (e.g. library listings).
  * Cache invalidation hooks in our repository layer ensure data freshness on writes.

  **PostgreSQL Partitioning & Indexing**:
  * Time‑series tables (expense logs) are automatically partitioned by month to keep queries sub‑second.
  * Composite indexes on (owner_id, created_at) speed up user‑scoped queries.
</details>

<details>
  <summary>🔒 Security</summary>

  **JWT Authentication**:
  * Auth0‑issued JWTs (RS256) with middleware handling validity of said JWTs.
  * Access and refresh tokens signed with RS256.
  * User‑context validation on every endpoint
  * DB‑backed user existence checks and route protection.

  **Input Validation & XSS Prevention**:
  * `bluemonday.StrictPolicy()` sanitizes any HTML inputs
  * Regex‑based whitelists for query parameters
  * Strict JSON schema and URL‑param validation

  **SQL Injection & Data Integrity**:
  * All DB interactions use parameterized queries and transactions
  * Foreign‑key constraints enforce referential integrity
  * User‑scoped filtering ensures zero lateral data access

  **Rate Limiting & Throttling**:
  * Redis‑backed rate limiter (configurable request per user)
  * Token‑bucket algorithm to smooth bursts

  **CORS & Network Protections**:
  * Fine‑grained CORS policy allows only approved origins
  * Private VPC connectivity (production) or Docker bridge (local)

  **Error Handling & Monitoring**:
  * Structured JSON logs with request IDs (ELK/Loki compatible)
  * Sentry.io integration for exception tracking and alerts

  **Cache Security**:
  * User‑scoped Redis keys with proper invalidation hooks
  * No sensitive data persisted in cache

</details>


<details>
  <summary>🚀 Infrastructure & DevOps</summary>

  **Observability & Monitoring**:
  * Structured logging with request tracing across frontend and backend
  * Custom log formats with user context, performance metrics, and error correlation
  * Centralized logging with Sentry integration for error tracking and alerting

  **CI/CD Pipeline**:
  * GitHub Actions with matrix builds (Go 1.19–1.21) and multi-arch Docker images
  * Database migrations and seeding with rollback capabilities
  * Backup/restore workflows with timestamped snapshots

  **Multi-Environment Containerization**:
  * Docker Compose orchestration with Traefik reverse proxy, PostgreSQL 14, Redis 7, and Sentry.io monitoring
  * Environment isolation with separate `.env` files for dev/test/prod
  * Health checks and graceful degradation across all services
  * Multi-stage Dockerfiles with optimized production builds

  **Cloudflare Deployments**:
  * Frontend & marketing on Pages with immutable previews
  * API gateway and edge functions on Workers for global low-latency
  * CDN optimization with automatic cache invalidation

  **Development Experience**:
  * Single-command workflows (`make dev-backend`, `make reset`)
  * Intelligent Vite control with network diagnostics and auto-recovery

</details>


## Documentation & References

- **Architecture Deep Dive**: [docs/architecture.md](./docs/architecture.md)
- **API Docs**: [docs/api/API-README.md](./docs/api/API-README.md)
- **OpenAPI Spec**: [docs/api/openapi.yaml](./docs/api/openapi.yaml)
- **Infrastructure**: [docs/api/containerized-env.md](./docs/containerized-env.md)
- **ADR Examples**: [docs/api/decisions/0001-initial-architecture.md](./docs/decisions/0001-initial-architecture.md)

## Setup & Installation

### Prerequisites
*  Docker (20.10.0+)
* Docker Compose (2.0.0+)
* Make

### Getting Started (Beta)
1. Pull down the full repo.

2. Check out the .env files for both `/backend` and `/frontend` directories.\
   QKO requires many env vars set up prior to attempting to build.\
   Please examine the `.env.sample` files to get started.

3. *Marketing + Frontend:*\
   After you have set up `.env` vars, check out the `vite-control.sh` script in both directories.\
   Use this build script to get you up and running quickly.\
   Both scripts come with complete documentation.\
   Run `./vite-control.sh help` in both directories to get started.

4. *Backend:*\
   Install Docker locally then look for the Makefile in the root directory.\
   This file also has complete documentation.\
   Run `make help` to get started.

5. Both Frontend + Backend must be running in order to run the build locally.

## Roadmap

*  ❤️ Wishlist Section - Save games to your wishlist and track when they go on sale across multiple vendor APIs
*  💹 More detailed expense reporting dashboard - download reports and track expenditures from all months of the year
*  ✨ Updated themes - Game inspired UI
*  🌐 Multi language support

## License

QKO is proprietary. All rights reserved.

---

**Created by [Ahn Ming Loke](https://github.com/loke)**  • [LinkedIn](https://www.linkedin.com/in/ahnmingloke/)

