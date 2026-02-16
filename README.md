# PG Management System - Backend

A high-performance microservice backend for managing Paying Guest (PG) facilities, built with Go and PostgreSQL.

## Authentication & Authorization

The system uses JWT (JSON Web Tokens) for secure API access.

### 1. Authentication Methods

#### **A. Google OAuth 2.0**
- **Endpoint**: `/auth/google/login`
- **Flow**: Redirects user to Google for login. Upon success, Google redirects back to `/auth/google/callback`, which returns a JWT token.

#### **B. Email & Password**
- **Registration**: `POST /auth/register` (Email, Password, Name)
- **Login**: `POST /auth/login` (Email, Password)
- **Response**: Returns a JWT token upon successful authentication.

### 2. Authorization (JWT)

All `/api/*` and `/graphql` routes are protected. To access them, include the token in your request header:

```http
Authorization: Bearer <your_jwt_token_here>
```

- **Token Expiry**: 24 hours.
- **Algorithm**: HS256.

## Project Structure
- `cmd/server`: Main entry point.
- `internal/database`: DB connection, schema, and repositories.
- `internal/handlers`: HTTP handlers (Auth, Rooms, Guests, Payments).
- `internal/middleware`: JWT authentication middleware.
- `internal/models`: Data structures.
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
   # Database
   DB_URL=postgres://user:password@localhost:5432/dbname
   PORT=8080

   # Authentication
   JWT_SECRET=your_super_secret_key
   
   # Google OAuth (Optional for local testing)
   GOOGLE_CLIENT_ID=your_id
   GOOGLE_CLIENT_SECRET=your_secret
   GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
   ```

2. **Run Server**:
   ```bash
   go run cmd/server/main.go
   ```

3. **API Examples**:
   - **Login**:
     ```bash
     curl -X POST http://localhost:8080/auth/login \
          -H "Content-Type: application/json" \
          -d '{"email": "user@example.com", "password": "password123"}'
     ```
   - **Get Rooms (Authenticated)**:
     ```bash
     curl http://localhost:8080/api/rooms \
          -H "Authorization: Bearer <TOKEN>"
     ```

## Version Control
This project uses Git for version control.
- `main`: Stable production-ready code.
- `develop`: Development branch for ongoing features and optimizations.
