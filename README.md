# PG Management System - Backend

A high-performance microservice backend for managing Paying Guest (PG) facilities, built with Go and PostgreSQL.

## Features
- **Room Management**: CRUD operations for PG rooms.
- **Guest Management**: Tracking guest details and room assignments.
- **Payment Tracking**: Managing payments per guest.
- **Performance Profiling**: Built-in `pprof` integration.
- **GraphQL & REST API**: Flexible data access layers.

## Tech Stack
- **Language**: Go (Golang)
- **Database**: PostgreSQL (Raw SQL with `pgx`)
- **Routing**: `gorilla/mux`
- **GraphQL**: `graphql-go`
- **Profiling**: `net/http/pprof`
- **CORS**: `rs/cors`

## Project Structure
- `cmd/server`: Main entry point.
- `internal/database`: DB connection and schema initialization.
- `internal/handlers`: HTTP request handlers.
- `internal/models`: Data structures.
- `internal/repository`: Data access layer.
- `internal/gql`: GraphQL schema and resolvers.
- `scripts`: Performance measurement and utility scripts.

## Optimization Techniques

### 1. Database Layer
- **Hand-written SQL**: Avoided ORM overhead by using raw SQL queries with `pgx`.
- **Schema Optimization**: Used appropriate data types and primary key indexing for fast lookups.
- **Automated Schema Initialization**: Ensures consistent database state on startup.

### 2. Service/API Layer
- **Concurrent Request Handling**: Leverages Go's goroutines and `net/http` native performance.
- **Efficient Routing**: Used `gorilla/mux` for structured and fast path matching.
- **Performance Monitoring**: Integrated `pprof` for real-time CPU and memory analysis.
- **Payload Minimization**: Structured JSON responses for minimal transfer overhead.

## Performance Metrics

Checked on: 2026-02-09

| Endpoint | Average Latency | Status |
|----------|-----------------|--------|
| `/health` | 4.31ms | PASS (< 200ms) |
| `/rooms` | 2.52ms | PASS (< 200ms) |
| `/guests` | 1.47ms | PASS (< 200ms) |

### CPU & Memory Profiling
- **CPU Profiling**: Integrated at `/debug/pprof/profile`.
- **Heap Analysis**: Available at `/debug/pprof/heap`.
- **Flame Graphs**: Can be generated via `go tool pprof`.

## Getting Started

1. **Environment Setup**:
   Create a `.env` file with:
   ```env
   DB_URL=postgres://user:password@localhost:5432/dbname
   PORT=8080
   ```

2. **Run Server**:
   ```bash
   go run cmd/server/main.go
   ```

3. **Measure Performance**:
   ```bash
   go run scripts/measure_performance.go
   ```

## Version Control
This project uses Git for version control.
- `main`: Stable production-ready code.
- `develop`: Development branch for ongoing features and optimizations.
