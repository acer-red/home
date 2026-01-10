# GitHub Copilot Instructions for Acer Official Engine

This repository contains the backend engine for the Acer Official project, built with Go, Gin, PostgreSQL, and Redis.

## Architecture Overview

- **Entry Point**: `engine/bin/server.go` initializes the application via `app.Main()`.
- **App Lifecycle**: `engine/internal/app` handles initialization of configuration, database, cache, and web server.
- **Web Layer**: `engine/service/web` contains HTTP handlers, middleware, and routing.
  - **Routes**: Defined in `route.go` files within subpackages (e.g., `engine/service/web/user/route.go`).
  - **Middleware**: `engine/service/web/common/middleware.go` handles cross-cutting concerns like CORS and Authentication.
- **Data Layer**: `engine/service/db` encapsulates PostgreSQL interactions (via GORM) and data models.
- **Cache Layer**: `engine/service/cache` manages Redis interactions, primarily for sessions.
- **Utilities**: `engine/util` provides shared helper functions and error definitions.

## Project Conventions

### API Response

- Use the `message` struct from `engine/service/web/common/msg.go` for all JSON responses.
- Standard format: `{ "err": <int>, "msg": <string>, "data": <interface> }`.
- Use helper functions like `msgOK()` or `msgCreated()` for success responses.

### Authentication

- Handled by `common.Auth()` middleware.
- Supports both **Cookie** (`jwt`) and **API Key** authentication.
- User context is stored in Gin context with key `"user"`.

### Database (PostgreSQL)

- Models are defined in `engine/service/db` using GORM (`gorm.Model`).
- Use `gorm` tags for database schema definition (e.g., `gorm:"uniqueIndex"`).
- `db.Init()` establishes the connection using `gorm.io/driver/postgres`.
- Database schema is auto-migrated via `db.AutoMigrate()`.
- Files are stored in the `files` table as `BYTEA` (replacing GridFS).
- Encapsulate DB logic within the `db` package methods.

### Configuration

- Config is defined in `engine/internal/app/config.go`.
- Loaded from a YAML file specified by the `-c` flag.
- `app.Config` struct mirrors the YAML structure.

### Logging

- Use `github.com/tengfei-xy/go-log`.
- Levels: `log.Debug`, `log.Info`, `log.Error`, `log.Fatal`.
- Log level is configurable via the `-v` flag.

## Integration Points

- **Redis**: Used for session management and caching. Connection initialized in `app.init_redis`.
- **PostgreSQL**: Primary data store. Connection initialized in `app.init_db`.

## Coding Guidelines

- **Error Handling**: Use errors defined in `engine/util/err.go` (e.g., `ErrNoFound`, `ErrInternalServer`) where applicable.
- **Context**: Pass `gin.Context` or `context.Background()` where needed, though GORM handles context implicitly in most cases.
- **Gin**: Use `gin.Context` for HTTP handlers. Avoid blocking operations in the main thread; use goroutines if necessary.
