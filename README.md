# LUXE — Backend API

Premium luxury eCommerce REST API built with **Go**, **Gin**, and **MongoDB Atlas**.

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.21+ |
| Framework | Gin |
| Database | MongoDB Atlas |
| Auth | JWT (access + refresh tokens) |
| Architecture | Clean Architecture (Handler → UseCase → Repository) |
| Validation | go-playground/validator |
| Logging | zerolog |
| API Docs | Swagger (swaggo) |
| Media | ImageKit CDN (+ local fallback) |
| Security | CORS, rate limiting, bcrypt passwords, RBAC |

## Project Structure

```
backend/
├── cmd/server/          # Application entry point
├── config/              # Environment configuration
├── internal/
│   ├── delivery/http/   # Gin handlers & router
│   ├── delivery/middleware/
│   ├── domain/          # Entities & repository interfaces
│   ├── usecase/         # Business logic
│   ├── repository/      # MongoDB implementations
│   └── pkg/             # JWT, imagekit, pagination, response
├── docs/                # Swagger generated docs
├── uploads/             # Local upload fallback (gitignored)
└── .env.example         # Environment template
```

## Features

- **Auth:** register, login, forgot/reset password
- **Storefront:** products, categories, cart, orders, reviews, coupons, banners
- **Admin:** dashboard, products, categories, orders, customers, coupons, banners, reviews
- **Uploads:** ImageKit CDN when configured; falls back to local `./uploads`
- **Middleware:** JWT auth, admin-only routes, rate limiter, CORS, request logging

## Prerequisites

- Go 1.21+
- MongoDB Atlas cluster (or local MongoDB)
- ImageKit account (optional, for CDN uploads)

## Setup

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your values (MongoDB URI, JWT secrets, ImageKit keys).

3. Install dependencies & run:
   ```bash
   go mod download
   go run ./cmd/server
   ```

4. API available at `http://localhost:8080`
   - Health: `GET /health`
   - Swagger: `GET /swagger/index.html`
   - Base API: `/api/v1`

## Environment Variables

See `.env.example` for all options. Required:

| Variable | Description |
|----------|-------------|
| `MONGO_URI` | MongoDB Atlas connection string |
| `MONGO_DB` | Database name (default: `luxe_db`) |
| `JWT_SECRET` | JWT signing secret (min 32 chars) |
| `JWT_REFRESH_SECRET` | Refresh token secret |
| `FRONTEND_URL` | Frontend origin for CORS |
| `IMAGEKIT_PRIVATE_KEY` | ImageKit private key (uploads) |
| `IMAGEKIT_PUBLIC_KEY` | ImageKit public key |
| `IMAGEKIT_URL_ENDPOINT` | ImageKit URL endpoint |

> **Never commit `.env` to Git.** It is listed in `.gitignore`.

## Docker

From project root:
```bash
docker compose up --build
```

## API Overview

| Group | Prefix | Auth |
|-------|--------|------|
| Auth | `/api/v1/auth` | Public |
| Products | `/api/v1/products` | Public |
| Categories | `/api/v1/categories` | Public |
| Cart / Orders | `/api/v1/cart`, `/orders` | JWT |
| Admin | `/api/v1/admin/*` | JWT + Admin role |

## Security Notes

- Rotate JWT secrets and database passwords if they were ever exposed.
- ImageKit private key must stay on the server only.
- Use `GIN_MODE=release` in production.

## License

MIT
