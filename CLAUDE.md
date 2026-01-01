# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI Proxy (cli-proxy-go) is a lightweight Go-based AI gateway service that provides unified API access to multiple AI platforms (Claude, OpenAI, Gemini, Bedrock, Azure) with API Key management and usage tracking.

**简化版本** - 移除了用户系统、套餐系统、邮箱验证等功能，专注于 AI 网关核心功能。

## Tech Stack

- **Backend**: Go 1.24, Gin 1.11, GORM (MySQL 8.0+), JWT, Zap logger
- **Frontend**: Vue 3.4 (Composition API), Vite 5, Element Plus, Alova, Pinia
- **Deployment**: Docker multi-stage build, Docker Compose

## Common Commands

### Backend

```bash
make run              # Run server directly
make dev              # Development mode with hot reload (requires 'air')
make build            # Build binary to bin/server
make test             # Run tests
make fmt              # Format code
make lint             # Run linter
go run ./cmd/server   # Direct Go run
```

### Frontend

```bash
cd web
npm install           # Install dependencies
npm run dev           # Dev server (port 3000, proxies to localhost:8080)
npm run build         # Production build
```

### Docker

```bash
docker-compose up -d                    # Start full stack (MySQL + app)
docker-compose up -d --build            # Rebuild and start
docker-compose logs -f cli-proxy        # View logs
```

## Architecture

### Directory Structure

```
cmd/server/main.go          # Application entry point
internal/
  ├── config/               # YAML configuration loading
  ├── handler/              # HTTP handlers (routes.go is the routing center)
  ├── middleware/           # JWT, API Key auth, CORS
  ├── model/                # GORM data models
  ├── service/              # Business logic layer
  ├── repository/           # Data access layer (MySQL)
  ├── proxy/
  │   ├── adapter/          # Platform adapters (Claude, OpenAI, Azure, Bedrock, Gemini)
  │   └── scheduler/        # Account selection, retry, token management
  └── cache/                # In-memory caching
pkg/
  ├── logger/               # Zap-based logging with rotation
  ├── response/             # Unified API response wrapper
  └── utils/                # JWT utilities
web/src/
  ├── views/                # Vue page components (admin only)
  ├── api/index.js          # API client (Alova)
  ├── router/index.js       # Vue Router configuration
  └── stores/               # Pinia state management
```

### Request Flow

1. Request enters through middleware chain: Logger → Recovery → CORS → Gzip
2. **Proxy routes** (`/claude/*`, `/openai/*`, `/gemini/*`): API Key auth → Client filter → ProxyHandler
3. **Admin routes** (`/api/*`): JWT auth → Handler

### Platform Adapters

Each AI platform has an adapter in `internal/proxy/adapter/`:
- `claude.go` - Claude Official/Console API
- `openai.go` - OpenAI Chat Completions
- `azure.go` - Azure OpenAI
- `bedrock.go` - AWS Bedrock
- `gemini.go` - Google Gemini
- `openai_responses.go` - OpenAI Responses API (for Codex CLI)

### Account Scheduler

`internal/proxy/scheduler/scheduler.go` handles:
- Account selection based on platform, status, and load
- Automatic retry with account failover
- Rate limit detection and account status updates
- OAuth token refresh management

### Authentication

- **JWT**: For admin login (`/api/*` routes), 24-hour expiry
- **API Key**: For proxy endpoints, stored in `x-api-key` header
- Admin credentials stored in config file (`configs/config.yaml`)

### Usage Tracking

- `DailyUsage`: Aggregated daily stats per API Key
- `UsageRecord`: Individual request records (persisted to MySQL)
- Pricing calculated per model with input/output token rates

## Key Files for Common Tasks

| Task | Files |
|------|-------|
| Add new API endpoint | `internal/handler/routes.go`, new handler in `internal/handler/` |
| Add new AI platform | `internal/proxy/adapter/`, update `internal/model/account.go` |
| Add new data model | `internal/model/`, `internal/repository/`, `internal/service/` |
| Add admin page | `web/src/views/`, `web/src/router/index.js` |
| Modify middleware | `internal/middleware/` |
| Change config options | `internal/config/config.go`, `configs/config.yaml` |

## Configuration

Main config file: `configs/config.yaml`

Key settings:
- `server.port`: HTTP port (default 8080)
- `mysql.*`: Database connection
- `jwt.secret`: JWT signing key
- `admin.username`: Admin username (default: admin)
- `admin.password`: Admin password (plain text or bcrypt hash)
- `cache.*`: TTL settings for sessions

Environment variables:
- `JWT_SECRET` - Override JWT secret
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database config
- `ADMIN_USERNAME`, `ADMIN_PASSWORD` - Admin credentials

## API Endpoints

### Proxy (requires API Key in `x-api-key` header)
- `POST /claude/v1/messages` - Claude API
- `POST /openai/v1/chat/completions` - OpenAI Chat
- `POST /v1/responses` - OpenAI Responses API
- `POST /gemini/v1/chat` - Gemini API

### Auth (public)
- `POST /api/auth/login` - Admin login

### Admin (JWT required)
- `/api/api-keys` - API Key management (CRUD)
- `/api/usage/*` - Usage statistics
- `/api/accounts` - AI platform account management
- `/api/models-config` - Model pricing configuration
- `/api/monitor` - System monitoring
- `/api/cache` - Cache management
- `/api/configs` - System configuration

## Default Credentials

- Username: `admin`
- Password: `admin123` (configurable via `admin.password` in config)

## Code Conventions

- Go code uses Chinese comments throughout
- Frontend uses Vue 3 Composition API with `<script setup>`
- API responses follow unified format via `pkg/response/`
- Database migrations handled automatically by GORM AutoMigrate

## Features

### Retained Features
- ✅ AI Gateway: Claude/OpenAI/Gemini/Bedrock/Azure proxy
- ✅ Multi-account rotation, OAuth auth, health checks
- ✅ Admin system: JWT auth (credentials in config file)
- ✅ Admin login captcha: Image captcha for admin login security
- ✅ API Key management: Admin creates and distributes keys
- ✅ API Key permissions: Platform/model/rate/client restrictions
- ✅ Usage statistics: Token count, request count, cost (per API Key)
- ✅ API Key quotas: Monthly/daily limits
- ✅ Admin dashboard (simplified)

### Removed Features
- ❌ User system: Registration, login, user management
- ❌ Package system: Subscription plans
- ❌ Email verification
- ❌ User portal (/user/*)
- ❌ Price multiplier system
