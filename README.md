# Signage Microservices Architecture

This repository contains two independent microservices built with Go, designed for high-scale digital signage deployments.

## Architecture

The project follows **Clean Architecture / Hexagonal Architecture** principles. Both services (Screen Service and Content Service) are separated by domain but reside in a monorepo for easier management.

- **Screen Service (`cmd/screen-service`)**: Exposes a public **HTTP API** (port `8001`) for external clients. It manages screen lifecycles and acts as a gateway for content assignment, internally communicating with the Content Service via gRPC.
- **Content Service (`cmd/content-service`)**: Manages media content and schedules via an internal **gRPC server** (port `50052`).

### API Documentation (Swagger)

The Screen Service includes integrated Swagger documentation.

- **Swagger UI**: Accessible at `http://localhost:8001/swagger/` when the service is running.
- **Update Docs**: To update the documentation after changing handler comments, run:
  ```bash
  make generate-swagger
  ```

### Key Endpoints (Screen Service)

| Method | Endpoint                | Description                                                               |
| ------ | ----------------------- | ------------------------------------------------------------------------- |
| `GET`  | `/health`               | Health Check                                                              |
| `POST` | `/screens`              | Create a new screen                                                       |
| `GET`  | `/screens/{id}`         | Get screen details                                                        |
| `PUT`  | `/screens/{id}/status`  | Update screen status (ONLINE/OFFLINE/MAINTENANCE)                         |
| `GET`  | `/screens/{id}/content` | Get all content assigned to a screen (via internal gRPC)                  |
| `POST` | `/screens/{id}/content` | **Composite**: Create new content AND assign it to the screen in one call |

## Running Locally

1. **Spin up a local Postgres instance**:
   ```bash
   docker run --name signage-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
   ```
2. **Generate All** (Protos & Swagger):
   ```bash
   make generate
   ```
3. **Run Screen Service** (HTTP :8001):
   ```bash
   make run-screen
   ```
4. **Run Content Service** (gRPC :50052):
   ```bash
   make run-content
   ```
5. **Run Tests**:
   ```bash
   make test
   ```

### Design Decisions & Best Practices

1. **Clean Architecture**: Deep logical separation per service under `internal/`. Each service contains `domain` (entities), `app` (business use cases), `infrastructure` (repositories), and `interfaces` (delivery layer).
2. **Hybrid Communication**: Public HTTP for external accessibility balanced with internal gRPC for performance and strict schema validation between services.
3. **Database Segregation**: Even though services share a Postgres instance, tables are strictly owned and accessed by their respective domains.
4. **Structured Logging (`zap`)**: High-performance structured logging.
5. **Config Management (`envconfig`)**: Twelve-factor app adherence with env-driven configurations.
6. **Graceful Shutdown**: Safe termination handling for containerized environments.
