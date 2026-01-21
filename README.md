# Slack Clone Backend

Production-ready Slack clone backend built with Golang, Gin, PostgreSQL, Redis, and WebSocket.

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- (Optional) Docker & Docker Compose

### Installation

```bash
# 1. Navigate to backend directory
cd backend

# 2. Copy environment file
cp .env.example .env

# 3. Edit .env with your database credentials
#    Update DB_USER, DB_PASSWORD, JWT_SECRET, etc.

# 4. Download dependencies
go mod download

# 5. Run database migrations
psql -U postgres -d slack_clone -f migrations/001_initial_schema.up.sql

# 6. Run the server
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

### Using Docker Compose

```bash
# From project root (slack/)
docker-compose up -d

# Check logs
docker-compose logs -f backend
```

## 📁 Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # ✅ Entry point
│
├── internal/
│   ├── config/
│   │   └── config.go               # ✅ Configuration
│   │
│   ├── database/
│   │   ├── postgres.go             # ✅ PostgreSQL connection
│   │   └── redis.go                # ✅ Redis connection
│   │
│   ├── middleware/
│   │   ├── auth.go                 # ⏳ TODO: JWT authentication
│   │   ├── cors.go                 # ⏳ TODO: CORS (using gin-cors now)
│   │   └── logger.go               # ⏳ TODO: Request logging
│   │
│   ├── models/
│   │   ├── models.go               # ✅ All data models
│   │   └── dto/                    # ⏳ TODO: Data Transfer Objects
│   │       ├── auth_dto.go
│   │       └── message_dto.go
│   │
│   ├── repository/
│   │   ├── user_repo.go            # ⏳ TODO: User database operations
│   │   ├── workspace_repo.go       # ⏳ TODO: Workspace operations
│   │   ├── channel_repo.go         # ⏳ TODO: Channel operations
│   │   └── message_repo.go         # ⏳ TODO: Message operations
│   │
│   ├── service/
│   │   ├── auth_service.go         # ⏳ TODO: Authentication logic
│   │   ├── user_service.go         # ⏳ TODO: User business logic
│   │   └── message_service.go      # ⏳ TODO: Message logic
│   │
│   ├── handler/
│   │   ├── auth_handler.go         # ⏳ TODO: Auth HTTP handlers
│   │   ├── user_handler.go         # ⏳ TODO: User handlers
│   │   ├── workspace_handler.go    # ⏳ TODO: Workspace handlers
│   │   └── message_handler.go      # ⏳ TODO: Message handlers
│   │
│   ├── websocket/
│   │   ├── hub.go                  # ⏳ TODO: WebSocket hub
│   │   ├── client.go               # ⏳ TODO: Client connection
│   │   └── handler.go              # ⏳ TODO: WS message handler
│   │
│   ├── webrtc/
│   │   └── signaling.go            # ⏳ TODO: WebRTC signaling
│   │
│   └── utils/
│       ├── response.go             # ⏳ TODO: Response helpers
│       └── errors.go               # ⏳ TODO: Error definitions
│
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                  # ✅ JWT utilities
│   │
│   └── hash/
│       └── hash.go                 # ✅ Password hashing
│
├── migrations/
│   ├── 001_initial_schema.up.sql   # ✅ Database schema
│   └── 001_initial_schema.down.sql # ✅ Rollback migration
│
├── go.mod                          # ✅ Dependencies
├── .env.example                    # ✅ Environment template
└── README.md                       # This file
```

## ✅ Completed Files

1. **`go.mod`** - All dependencies defined
2. **`.env.example`** - Environment configuration template
3. **`migrations/001_initial_schema.up.sql`** - Complete database schema
4. **`internal/config/config.go`** - Configuration management
5. **`internal/models/models.go`** - All data models
6. **`internal/database/postgres.go`** - PostgreSQL connection
7. **`internal/database/redis.go`** - Redis connection
8. **`pkg/jwt/jwt.go`** - JWT token management
9. **`pkg/hash/hash.go`** - Password hashing
10. **`cmd/server/main.go`** - Main server with route structure

## ⏳ TODO: Implementation Guide

### Phase 1: Authentication System

Create these files in order:

#### 1. DTOs (Data Transfer Objects)
```go
// internal/models/dto/auth_dto.go
package dto

type RegisterRequest struct {
    Email           string `json:"email" binding:"required,email"`
    Username        string `json:"username" binding:"required,min=3,max=50"`
    Password        string `json:"password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
    FullName        string `json:"full_name"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    User         *models.User `json:"user"`
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
}
```

#### 2. User Repository
```go
// internal/repository/user_repo.go
package repository

// Create, FindByEmail, FindByID, Update methods
// See BACKEND_GUIDE.md for detailed implementation
```

#### 3. Auth Service
```go
// internal/service/auth_service.go
package service

// Register, Login, RefreshToken, Logout methods
// Business logic layer
```

#### 4. Auth Handler
```go
// internal/handler/auth_handler.go
package handler

// HTTP handlers for auth endpoints
// Connects Gin routes to service layer
```

#### 5. Auth Middleware
```go
// internal/middleware/auth.go
package middleware

// JWT authentication middleware
// Validates token and sets user in context
```

### Phase 2: WebSocket Real-time

#### 1. WebSocket Hub
```go
// internal/websocket/hub.go
package websocket

// Hub manages all WebSocket connections
// Rooms for channels/DMs
// Broadcast messages to rooms
```

#### 2. WebSocket Client
```go
// internal/websocket/client.go
package websocket

// Client connection
// readPump and writePump goroutines
```

#### 3. WebSocket Handler
```go
// internal/websocket/handler.go
package websocket

// Handle WebSocket upgrade
// Process incoming messages
```

### Phase 3: Core Features

Implement in this order:
1. **Workspaces** - Repository, Service, Handler
2. **Channels** - Repository, Service, Handler
3. **Messages** - Repository, Service, Handler
4. **Reactions** - Repository, Service, Handler
5. **Attachments** - File upload, S3/MinIO integration

### Phase 4: WebRTC

```go
// internal/webrtc/signaling.go
package webrtc

// WebRTC signaling server
// Handle offer, answer, ICE candidates
```

## 🔧 Development Commands

```bash
# Run server
go run cmd/server/main.go

# Run with auto-reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server

# Format code
go fmt ./...

# Lint (install golangci-lint first)
golangci-lint run
```

## 🔐 Environment Variables

See `.env.example` for all available options.

**Important:**
- Change `JWT_SECRET` to a strong random string in production
- Use proper database credentials
- Update `ALLOWED_ORIGINS` for frontend URL

## 📚 API Documentation

Once handlers are implemented, API endpoints will be:

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout

### Users
- `GET /api/users/me` - Get current user
- `PUT /api/users/me` - Update profile

### Workspaces
- `GET /api/workspaces` - List user's workspaces
- `POST /api/workspaces` - Create workspace
- `GET /api/workspaces/:id` - Get workspace details

### Channels
- `GET /api/workspaces/:id/channels` - List channels
- `POST /api/workspaces/:id/channels` - Create channel
- `GET /api/channels/:id/messages` - Get messages
- `POST /api/channels/:id/messages` - Send message

### WebSocket
- `WS /ws` - WebSocket connection

### WebRTC
- `WS /webrtc/signaling` - WebRTC signaling

## 🐳 Docker

```bash
# Build image
docker build -t slack-backend-clone .

# Run container
docker run -p 8080:8080 --env-file .env slack-backend-clone
```

## 📖 Next Steps

1. Read **`BACKEND_GUIDE.md`** for detailed architecture and implementation patterns
2. Implement Phase 1 (Authentication) following the guide
3. Test authentication endpoints with Postman/Insomnia
4. Implement Phase 2 (WebSocket) for real-time messaging
5. Continue with Phase 3 (Core Features)
6. Add WebRTC support (Phase 4)

## 🤝 Contributing

This is a learning project. Feel free to:
- Add missing features
- Improve code structure
- Add tests
- Optimize performance

## 📝 License

MIT

---

**For detailed implementation guide, see:**
- `BACKEND_GUIDE.md` - Complete backend architecture and patterns
- `FRONTEND_GUIDE.md` - Frontend implementation guide
