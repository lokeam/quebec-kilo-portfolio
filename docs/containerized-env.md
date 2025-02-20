# Q-KO Docker Environment Documentation & Troubleshooting Guide

## Introduction

This document provides an overview of the containerized environment for Q-KO. It covers:
- The Docker Compose setup, including service definitions and network configurations.
- The structure and usage of the Dockerfiles and Makefile for local development and deployment.
- Steps for validating the configuration.
- A troubleshooting runbook detailing common issues—such as Traefik network misconfiguration—and their resolutions.

---

## Table of Contents

- [Overview](#overview)
- [Environment Setup](#environment-setup)
  - [Docker Compose Configuration](#docker-compose-configuration)
  - [Dockerfile Details](#dockerfile-details)
  - [Makefile Operations](#makefile-operations)
- [Network Configuration & Verification](#network-configuration--verification)
- [Validating Your Environment](#validating-your-environment)
- [Troubleshooting Runbook](#troubleshooting-runbook)
  - [Traefik Network Issues](#traefik-network-issues)
  - [General Troubleshooting Commands](#general-troubleshooting-commands)
- [Additional Resources](#additional-resources)

---

## Overview

Q-KO is a collector's companion designed for managing media collections (currently specializing in video games). Our containerized architecture includes:
- A **Traefik** reverse proxy for routing external traffic.
- A **Golang API** for backend operations.
- A **React Frontend** for user interaction.
- Supporting services like **Postgres**, **Redis**, and **Mailhog**.

This environment is defined in the `docker-compose.yml` file, which has been hardened to use explicit network names and network aliases (e.g., using the network `quebec-kilo_web` rather than a generic `web`).

> **Note:**
> Run `docker network ls` to list out the spread of networks for qko. Look for the exact names of the networks attached to your API container. For example, verify that the `web` network is actually named `quebec-kilo_web`.

---

## Environment Setup

### Docker Compose Configuration

The `docker-compose.yml` file defines both services and networks explicitly. Key points include:
- **Explicit Network Names:**
  We define networks with the `name` property to avoid Docker Compose project-prefix issues.
- **Network Aliases:**
  Services (like the API) use aliases to ensure consistent internal DNS resolution.
- **Traefik Labels:**
  Traefik is guided to use the proper network via the label `"traefik.docker.network=quebec-kilo_web"`.

Refer to our [Docker Compose file](../docker-compose.yml) for full details.

### Dockerfile Details

Each service has its own `Dockerfile` configured to build lightweight images. These files are designed to:
- Optimize build layers.
- Include app-specific health checks if needed.
- Provide entrypoints that integrate with our container orchestration.

### Makefile Operations

The Makefile (see [Makefile](../Makefile)) provides the following targets:
- `make up` and `make restart` — for building and launching backend services.
- `make dev`, `make test`, and `make prod` — to set up various environments.
- Health-check targets (e.g., `make health`, `make health-detail`, `make health-api`) that help in verifying that the containers are running as expected.

For more details, view the Makefile help by running `make help`.

---

## Network Configuration & Verification

### Expected Network Names

Our networks are defined explicitly:
- **Web Network:** `quebec-kilo_web`
- **Backend Network:** `quebec-kilo_backend`

### How to Verify

Before starting your services, confirm that the networks exist and have the expected names:
```bash
docker network ls
```

You should see entries similar to:
```bash
NETWORK ID NAME DRIVER SCOPE

0ee97acd7797 quebec-kilo_web bridge local
51159f4a4ed4 quebec-kilo_backend bridge local
```

If the names don’t match, review your `docker-compose.yml` and ensure there are no typos or inconsistencies.


## Validating Your Environment

Before deploying, run:
```bash
bash
docker-compose config
```

This command outputs the fully expanded configuration, allowing you to verify that your environment variables, network names, and volumes are set as expected.
**Tip:** Integrate this check in your CI/CD pipeline or via a Makefile target (e.g., `make validate`).


## Troubleshooting Runbook

In this section, we detail common issues and resolution steps. This runbook is an evolving document—update it as new issues arise.

### Traefik Network Issues

#### **Symptoms:**
- Traefik logs include warnings like:
  ```
  Could not find network named 'web' for container '/quebec-kilo-api-1'! Maybe you're missing the project's prefix in the label? Defaulting to first available network.
  ```
- HTTP requests to endpoints routed via Traefik (e.g., `api.localhost`) return 504 errors.

#### **Likely Causes:**
- The label `"traefik.docker.network=web"` was used instead of the explicit `"traefik.docker.network=quebec-kilo_web"`.
- Docker Compose automatically prefixes network names with the project name, so a network declared as `web` might appear as `quebec-kilo_web`.
  > **Remember:** Run `docker network ls` to verify network names.

#### **Resolution Steps:**
1. **Update the Label:**
   Ensure the API service label reads:
   ```yaml
   - "traefik.docker.network=quebec-kilo_web"
   ```
2. **Recreate Containers:**
   Run:
   ```bash
   docker compose down
   docker compose up --build -d
   ```
3. **Verify Networks:**
   Run `docker network ls` and ensure that `quebec-kilo_web` exists.
4. **Monitor Logs:**
   Check the Traefik container logs for any remaining network-related warnings.

---

### General Troubleshooting Commands

- **View Container Status:**
  ```bash
  docker compose ps
  ```
- **View Logs:**
  ```bash
  docker compose logs
  ```
- **Health Check:**
  ```bash
  make health
  make health-detail
  ```
- **Configuration Validation:**
  ```bash
  docker-compose config
  ```

---

## Additional Resources

- Docker Documentation: [Docker Compose](https://docs.docker.com/compose/)
- Traefik Documentation: [Traefik User Guide](https://doc.traefik.io/traefik/)
- CI/CD Tips for Docker: [Docker in CI/CD](https://docs.docker.com/ci-cd/)

---

## Conclusion

This document should serve as the single source of truth for developers working with the Q-KO containerized environment. It includes:
- A detailed explanation of how our environment is structured.
- Steps to validate configuration.
- A troubleshooting runbook for common issues like network misconfigurations (e.g., Traefik problems).
