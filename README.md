# EventFlow

Event ingestion and metrics API built with Go. Ingests events with idempotency, stores them in PostgreSQL and exposes aggregated metrics with optional caching and load balancing.

## Quick Start (How to Run)

| Mode | Command | Access |
|------|---------|--------|
| **Local** | `go run ./cmd/server` | http://localhost:8080 |
| **Docker (single app)** | `docker compose -f docker-compose.simple.yml up --build` | http://localhost:8080 |
| **Full stack (gateway + 3 replicas)** | `docker compose up --build` | Gateway: http://localhost:8080 |

Requires PostgreSQL. Copy `.env.example` to `.env` and set `DATABASE_URL`.

## Features

- **Event ingestion** — POST events with validation and idempotency (by `event_id` or `user_id|event_name|timestamp`)
- **Metrics API** — GET aggregated counts with optional `group_by` (e.g. channel)
- **CQRS** — Commands (write) and queries (read) separated
- **Caching** — In-memory TTL cache for metrics responses
- **Load balancer** — Round-robin gateway across multiple backend instances
- **Docker** — Full stack with PostgreSQL, app replicas and gateway

## Architecture of the event flow 

```
                    ┌─────────────┐
                    │   Gateway   │  :8080 (round-robin)
                    │ (optional)  │
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         ▼                 ▼                 ▼
   ┌──────────┐      ┌──────────┐      ┌──────────┐
   │  App 1   │      │  App 2   │      │  App 3   │  :8080 each
   └────┬─────┘      └────┬─────┘      └────┬─────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                           ▼
                    ┌─────────────┐
                    │  PostgreSQL │  :5432
                    └─────────────┘
```

# The gateway in not required you could build just the simple.yml


**Layers:**
- **WebAPI** — HTTP handlers, router and DTOS
- **Middleware** — Caching, Logs 
- **Application** — Commands (create_event) and queries (get_metrics)
- **Infra** — Domain, contracts, repositories, migrations and config data

## Tech Stack

- **Go 1.24** — Echo, GORM, pgx
- **PostgreSQL 16** — Persistent storage
- **GORM** — Model-based CRUD and simple queries
- **Raw SQL** — Custom aggregations (e.g. `COUNT(DISTINCT)`, `COALESCE`)

## Project Structure

```
cmd/
├── server/     # Main API server
└── gateway/     

webapi/         # HTTP layer
├── handlers/    
└── router.go

application/
├── commands/create_event/   # Write path
└── queries/get_metrics/     # Read path

infra/
├── domain/     # Event model
├── contracts/  # Repository interfaces
├── persistence/repositories/
├── migrations/
└── config.go

middleware/cache/   # In-memory TTL cache
tests/
├── unit/       # Service unit tests
└── load_tests/ # Load tests
        ____

## Prerequisites

- Go 1.24+
- PostgreSQL 16+ (or use Docker)
- Docker & Docker Compose (optional)

## Setup

### 1. Environment

```bash
cp .env.example .env
```

Edit `.env` with your values. Key variables:

| Variable | Description | Default |
|---------|-------------|---------|
| `HTTP_PORT` | API server port | 8080 |
| `DATABASE_URL` | PostgreSQL connection string | — |
| `CACHE_TTL_SECONDS` | Metrics cache TTL | default  60 |
| `GATEWAY_PORT` | Load balancer port | 8080 |
| `BACKEND_URLS` | Comma-separated backend URLs | — |

### 2. Run locally (single app)

Ensure PostgreSQL is running, then:

```bash
go run ./cmd/server
```

### 3. Run with Docker (single app)

```bash
docker compose -f docker-compose.simple.yml up --build
```

App: `http://localhost:8080`

### 4. Run full stack (gateway + 3 app replicas)

```bash
docker compose up --build
```

- Gateway: `http://localhost:8080`
- Backends: app1, app2, app3 (internal)

## API Reference

### Health

```http
GET /health
```

Returns 200 if DB is reachable.

---

### Ingest event

```http
POST /events
Content-Type: application/json
```

**Request body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `event_id` | string | No | Idempotency key (if provided) |
| `user_id` | string | Yes | User identifier |
| `event_name` | string | Yes | Event type (e.g. `click`, `purchase`) |
| `timestamp` | int64 | Yes | Unix timestamp |
| `channel` | string | No | Platform (e.g. `web`, `mobile`) |
| `campaign_id` | string | No | Campaign identifier |
| `tags` | []string | No | Tags |
| `metadata` | object | No | Extra key-value data |

**Example:**

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "event_name": "click",
    "timestamp": 1771545600,
    "channel": "web"
  }'
```

**Response:** `202 Accepted` with `{"status":"accepted"}`

**Idempotency:** Duplicate events (same `event_id` or same `user_id|event_name|timestamp`) return `409 Conflict`.

---

### Get metrics

```http
GET /metrics?event_name=click&from=1708400000&to=1708600000&group_by=channel
```

**Query params:**

| Param | Type | Required | Description |
|-------|------|----------|-------------|
| `event_name` | string | Yes | Event type to aggregate |
| `from` | int64 | Yes | Start timestamp (inclusive) |
| `to` | int64 | Yes | End timestamp (inclusive) |
| `group_by` | string | No | Group by field (e.g. `channel`) |

**Example:**

```bash
curl "http://localhost:8080/metrics?event_name=click&from=1771545600&to=1771545800&group_by=channel"
```

**Response:**

```json
{
  "total_count": 150,
  "unique_user_count": 42,
  "group_by": "channel",
  "buckets": [
    {"key": "web", "count": 90}
  ]
}
```

## Tests

```bash
# Unit tests
go test ./tests/unit/... -v -short

# Load tests (longer)
go test ./tests/load_tests/... -v
```

## Assumptions, Design Decisions 

### Assumptions

- **Idempotency**: Clients may retry; duplicates are identified by `event_id` or `user_id|event_name|timestamp`.
- **Timestamps**: Unix seconds; client-provided. No server-side clock correction.
- **Scale**: Single PostgreSQL instance; in-memory cache per app instance (no distributed cache).
- **Metrics**: Aggregations are computed on read; no pre-aggregation or materialized views.

### Design Decisions

- **CQRS**: Commands (create_event) and queries (get_metrics) separated for clear read/write boundaries.
- **DTOs in application layer**: Request/response DTOs live in commands/queries; handlers bind HTTP and delegate to services.
- **GORM + raw SQL**: GORM for CRUD; raw SQL for aggregations (`COUNT(DISTINCT)`, `GROUP BY`).
- **Round-robin gateway**: Simple load balancing; no health checks or weighted routing.
- **Cache**: Metrics cache is in the get_metrics service, not HTTP middleware. 
---

## Example Requests (curl)

**Ingest event:**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user-123","event_name":"click","timestamp":1771545600,"channel":"web"}'
```

**Get metrics:**
```bash
curl "http://localhost:8080/metrics?event_name=click&from=1771545600&to=1771545800&group_by=channel"
```

**Postman**: Import `postman/EventFlow.postman_collection.json` for ready-made requests.

---

## Roadmap (TODO Items)

- **Platform filtering** — Add `channel` query param to GET /metrics to filter by platform
- **Batch ingestion** — POST /events/batch for higher throughput
- **Cache invalidation** — Invalidate metrics cache when a new event is created
- **Rate limiting** — Add rate limiting middleware
- **In-memory cache**: Fast but not shared across replicas; each app has its own cache.  
 

## Trade-offs Under Time Pressure

- **No rate limiting**: I have to skipped to focus on core features. 
- **No batch ingestion**: Single event api only batching would be improve throughput  
- **Limited tests**: More unit tests (services or handlers ) .
- **No API versioning**: Versioning Endpoints would add prefix for prod 
- **Gateway has no health checks**: Doesnt removing unhelathy backends form rotation

## Alternative Approaches Considered

- **Event sourcing / Kafka**: Chosed PostgreSQL for simplicity and event streaming would suit.
- **Redis for cache**: Redis would enable shared cache. 
- **Nginx/Envoy as gateway**: Nginx or others would be best chooise for pure go
 

## Production-Grade Improvements

- **Observability**: Structured logging, metrics for prometheus and distributed tracing.
- **Security**: API keys, JWT, rate limiting, input sanitization, CORS. 
- **Deployment**: Adding Kubernetes manifests and CI/CD pipeline. 
- **Testing**: Integration tests, contract tests, chaos testing, load test automation in CI.
- **CodeS Security**: Using Jenkis etc.


