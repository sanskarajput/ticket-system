# Ticket System

A REST API service for a ticket management system built in Go with SQLite.

## Stack

- **Language:** Go 1.26
- **Database:** SQLite (via `go-sqlite3`)
- **Auth:** JWT (HS256, 24h expiry)
- **Password storage:** HMAC-SHA256 with random salt

## Local Run

```bash
go run ./cmd/main.go
```

## Docker Run

```bash
docker build -t ticket-system .
docker run -p 8080:8080 ticket-system

# or just
docker compose up --build
```

Health check:
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Environment Variables

| Variable     | Default                          | Description             |
|--------------|----------------------------------|-------------------------|
| `PORT`       | `8080`                           | Server port             |
| `DB_PATH`    | `./tickets.db`                   | SQLite database path    |
| `JWT_SECRET` | `random-secret-key`| JWT signing secret      |

Copy `.env.example` to `.env` and fill in values for production.

## API Endpoints

### Health

```
GET /health
```

### Auth

```
POST /auth/register
Body: { "email": "user@example.com", "password": "secret123" }

POST /auth/login
Body: { "email": "user@example.com", "password": "secret123" }
Response: { "token": "<jwt>" }
```

### Tickets (require `Authorization: Bearer <token>`)

```
POST   /tickets              Create ticket
GET    /tickets              List my tickets
GET    /tickets/{id}         Get my ticket by ID
PATCH  /tickets/{id}/status  Update ticket status
```

#### Create ticket

```json
{ "title": "Bug in login", "description": "Login fails on Safari" }
```

#### Update status

```json
{ "status": "in_progress" }
```

**Status flow:** `open` → `in_progress` → `closed` (one-way, no reversal)


---
