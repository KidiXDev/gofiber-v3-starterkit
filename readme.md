# GoFiber v3 Starterpack

A production-ready Go Fiber v3 boilerplate with PostgreSQL, Redis, S3/MinIO, JWT authentication with Redis token management, and clean architecture.

## Features

- **Go Fiber v3** - Fast HTTP framework
- **PostgreSQL + Bun ORM** - Database with powerful ORM
- **Redis** - Token management and caching
- **S3/MinIO** - Object storage for file uploads
- **JWT Authentication** - Access & Refresh tokens stored in Redis
- **Dependency Injection** - Using `go.uber.org/dig`
- **Atlas Migrations** - Database schema management
- **Zerolog** - Structured logging with colorful console output
- **Graceful Shutdown** - Proper signal handling

## Project Structure

```
.
├── main.go                          # Application entry point
├── go.mod                           # Go module dependencies
├── atlas.hcl                        # Atlas migration configuration
├── .env.example                     # Environment template
├── rename-module.sh                 # Linux/macOS module rename script
├── rename-module.bat                # Windows module rename script
├── app/
│   ├── api/
│   │   ├── controllers/             # HTTP request handlers
│   │   ├── services/                # Business logic layer
│   │   └── types/                   # Request/Response DTOs
│   ├── models/                      # Database models (Bun ORM)
│   ├── routes/                      # Route registration
│   └── shared/                      # Shared utilities (responses, errors)
├── pkg/
│   ├── client/
│   │   ├── db/                      # PostgreSQL client
│   │   ├── redis/                   # Redis client
│   │   └── s3/                      # S3/MinIO client
│   ├── config/                      # Fiber & CORS configuration
│   ├── middlewares/                 # Auth, validation middlewares
│   └── utils/                       # Utility functions
├── migrations/                      # SQL migration files
│   ├── migrate.go                   # Migration runner
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
├── loader/                          # Atlas schema loader
│   └── main.go
└── hc/                              # Health check utility
    └── main.go
```

## Quick Start

### ⚡ One-Line Installation

**Linux / macOS:**

```bash
/bin/bash <(curl -fsSL https://raw.githubusercontent.com/KidiXDev/gofiber-v3-starterkit/main/install.sh)
```

**Windows (PowerShell):**

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/KidiXDev/gofiber-v3-starterkit/main/install.ps1'))
```

This will:
1. Clone the repository
2. Ask for your project name
3. **Automatically rename the module** to your desired path
4. Prepare the environment

---

### Manual Setup

If you prefer to clone manually:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/KidiXDev/gofiber-v3-starterkit.git my-app
   cd my-app
   rm -rf .git  # Optional: Start with fresh git history
   git init
   ```

2. **Rename the Module:**
   
   **Linux/macOS:**
   ```bash
   chmod +x rename-module.sh
   ./rename-module.sh github.com/yourusername/your-project
   ```
   
   **Windows:**
   ```batch
   rename-module.bat github.com/yourusername/your-project
   ```

3. **Setup Environment:**

   ```bash
   cp .env.example .env
   ```

Edit `.env` with your configuration:

```env
# App
APP_ENV=development
APP_HOST=127.0.0.1
APP_PORT=8000

# PostgreSQL
DATABASE_URL=postgres://user:password@localhost:5432/mydb?sslmode=disable
DEV_DATABASE_URL=postgres://user:password@localhost:5432/mydb_dev?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# S3/MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET_NAME=uploads
MINIO_SECURE=false

# JWT
JWT_SECRET=your-super-secret-key-change-this
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run Database Migrations

**Option A: Using Atlas CLI (Recommended)**
```bash
# Install Atlas
curl -sSf https://atlasgo.sh | sh

# Apply migrations
atlas migrate apply --env local
```

**Option B: Using the bun ORM migration runner**
```bash
go run migrations/migrate.go
```

**With seeder:**
```bash
go run migrations/migrate.go --seed
```


### 6. Run the Application

**Development:**
```bash
go run .
```

**Development (with Hot Reload):**

1. Install [Air](https://github.com/cosmtrek/air):
```bash
go install github.com/cosmtrek/air@latest
```

2. Run with Air:
```bash
air
```

**Build & Run:**
```bash
go build -o app .
./app
```

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/livez` | Health check |

### Protected Endpoints (Require Bearer Token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/me` | Get current user |
| PUT | `/api/v1/auth/me` | Update profile |
| POST | `/api/v1/auth/logout` | Logout (revoke token) |
| POST | `/api/v1/auth/logout-all` | Logout from all devices |
| GET | `/api/v1/users` | List all users |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

## Authentication Flow

### Token Storage in Redis

```
┌─────────────────────────────────────────────────────────┐
│                    Redis Keys                           │
├─────────────────────────────────────────────────────────┤
│ access_token:{tokenID}  → userID   (TTL: 15 min)       │
│ refresh_token:{tokenID} → userID   (TTL: 7 days)       │
│ user_refresh:{userID}   → Set of refresh token IDs     │
└─────────────────────────────────────────────────────────┘
```

### Login Response

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe"
    },
    "auth": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
    }
  }
}
```

### Making Authenticated Requests

```bash
curl -H "Authorization: Bearer <access_token>" \
     http://localhost:8000/api/v1/auth/me
```

### Refresh Token

```bash
curl -X POST http://localhost:8000/api/v1/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<refresh_token>"}'
```

## Atlas Migrations

### Generate New Migration

```bash
atlas migrate diff migration_name --env bun
```

### Apply Migrations

```bash
atlas migrate apply --env local
```

### Check Migration Status

```bash
atlas migrate status --env local
```

## Adding New Models

1. Create model in `app/models/`
2. Add model to `loader/main.go`
3. Generate migration: `atlas migrate diff add_new_model --env bun`
4. Apply: `atlas migrate apply --env local`

## Health Check

Build and use the health check utility for Docker or monitoring:

```bash
go build -o hc ./hc
./hc                           # Check default localhost:8000/livez
./hc http://myserver:8000/livez  # Check custom URL
```
