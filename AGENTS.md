# AGENTS.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Kanban board application for self-hosted GitLab with real-time synchronization. The project consists of:
- **Backend**: Go-based REST API (Gin framework) with background worker for GitLab sync
- **Frontend**: Angular 19 single-page application
- **Database**: Supports SQLite, PostgreSQL, or MySQL via GORM

## Development Commands

### Backend (Go)

```bash
# Run backend in development mode
./scripts/run_dev_backend.sh

# Or manually from backend directory
cd backend
go run main.go

# Build backend
cd backend
go build -o board main.go
```

### Frontend (Angular)

```bash
# Run frontend in development mode
./scripts/run_dev_frontend.sh

# Or manually from frontend/board directory
cd frontend/board
npm install
npm run start        # Dev server on http://localhost:4200
npm run build        # Production build
npm run watch        # Build with watch mode
npm test            # Run Jasmine/Karma tests
```

### Docker

```bash
# Run with SQLite (default)
docker-compose up -d

# Run with PostgreSQL
docker-compose -f docker-compose.postgresql.yml up -d

# Run with MySQL
docker-compose -f docker-compose.mysql.yml up -d
```

## Environment Configuration

Copy `.env.example` to `.env` and configure:

**Required settings:**
- `GITLAB_URL`: Your GitLab instance URL (e.g., https://gitlab.yourdomain.com)
- `GITLAB_TOKEN`: GitLab private token for API access
- `API_SECRET`: Secret key for JWT token generation

**Key optional settings:**
- `DB_TYPE`: Database type (sqlite, postgresql, mysql)
- `GITLAB_SYNC_ENABLED`: Enable/disable background sync (default: true)
- `GITLAB_SYNC_PERIOD_MIN`: Sync interval in minutes (default: 10)
- `ALLOW_ORIGINS`: CORS origins (comma-separated)

## Architecture

### Backend Architecture (Clean Architecture Pattern)

The backend follows a layered architecture with dependency injection using `goioc/di`:

**Layers:**
1. **`internal/interfaces/api`**: HTTP layer (Controllers, DTOs, Middleware)
   - Controllers handle HTTP requests and call use cases
   - JWT middleware for authentication
   - CORS configuration

2. **`internal/app`**: Application layer (Use Cases, Services)
   - Business logic orchestration
   - Use cases coordinate between domain and infrastructure

3. **`internal/domain`**: Domain layer (Models, Business Rules)
   - Core domain entities and business logic

4. **`internal/infra`**: Infrastructure layer
   - `db/`: Database connection management (SQLite, PostgreSQL, MySQL)
   - `persistance/`: Repositories and database models (GORM)
   - `persistance/spec/`: Query specifications for complex queries
   - `gitlab/`: GitLab API client integration
   - `cache/`: In-memory caching layer
   - `services/`: Infrastructure services (JWT, Password hashing)

5. **`cmd/`**: Application entry points
   - `api.go`: REST API server (Gin)
   - `worker.go`: Background worker for GitLab synchronization

**Key patterns:**
- Dependency injection configured in `main.go`
- Repository pattern for data access
- Specification pattern for complex queries (see `internal/infra/persistance/spec/`)
- Use cases encapsulate business operations

### Frontend Architecture (Angular)

**Module structure:**
- **`modules/core`**: Core shared services and components
  - `services/`: Base services (JWT, Authentication, HTTP base, Toast notifications)
  - `interceptors/`: HTTP interceptors for auth headers
  - `models/`: Core domain models

- **`modules/kanban`**: Main kanban board feature module
  - `components/`: 28+ kanban-specific components (cards, columns, modals, lists)
  - `services/`: Domain services (Epic, Issue, Sprint, Team, User, Label, etc.)
  - `pipes/`: Custom Angular pipes
  - `models/`: Kanban domain models

- **`modules/ui`**: Reusable UI components (buttons, inputs, modals)
- **`modules/bootstrap-ui`**: Bootstrap-based UI wrapper components

**Key services:**
- `BaseService`: Generic HTTP service with caching (`CollectionCache`)
- `AccountService`: Authentication and user session management
- `JwtService`: JWT token handling and storage
- Domain-specific services extend `BaseService` for CRUD operations

**Routing:**
- `/` - Main kanban board
- `/auth` - Login page
- `/sprints` - Sprint management
- `/reports` - Reports view

### Data Synchronization

The backend runs a background worker (`cmd/worker.go`) that:
- Syncs with GitLab API on a configurable interval
- Updates local database with GitLab issues, users, projects, labels, etc.
- Operates independently from the API server
- Uses `SyncUseCases` to orchestrate the sync process

### Database Schema

Key entities:
- **Issue**: GitLab issues with kanban metadata
- **Epic**: Issue epics for grouping
- **Sprint**: Time-boxed iterations
- **Column**: Kanban board columns
- **Label**: GitLab labels
- **User**: GitLab users
- **Team**: User teams
- **Project**: GitLab projects
- **Group**: GitLab groups
- **Release**: GitLab releases
- **IssueBinding**: Relationships between issues and sprints/epics

All models are in `backend/internal/infra/persistance/models/`

## Making Changes

### Adding a New Backend Feature

1. Define the domain model in `internal/domain/models/`
2. Create database model in `internal/infra/persistance/models/`
3. Create repository in `internal/infra/persistance/repo/`
4. Add query specification in `internal/infra/persistance/spec/` if complex queries needed
5. Create use case in `internal/app/usecases/`
6. Create controller and DTOs in `internal/interfaces/api/`
7. Register all beans in `main.go` (repositories, queries, use cases, controllers)
8. Add controller to router initialization in `cmd/api.go`

### Adding a New Frontend Feature

1. Create service in appropriate module's `services/` directory
2. Create models in module's `models/` directory
3. Create components using Angular CLI or manually
4. Update module imports in the feature module file (e.g., `kanban.module.ts`)
5. Add routes in `app-routing.module.ts` if needed

### Dependency Injection

The backend uses `goioc/di`. Register beans in `main.go`:
```go
di.RegisterBean("BeanName", reflect.TypeOf((*YourType)(nil)))
di.RegisterBeanInstance("instanceName", instance)
di.RegisterBeanFactory("factoryName", di.Singleton, factoryFunc)
```

Access beans:
```go
bean := di.GetInstance("BeanName").(*YourType)
```

## Testing

- Frontend tests use Jasmine/Karma framework
- Run frontend tests: `cd frontend/board && npm test`
- Backend has no test files currently - add `_test.go` files as needed

## Notes

- The backend uses GORM for ORM with auto-migration
- Frontend environment configuration is injected at runtime via `window['env']`
- The application supports dark/light theme switching
- Built-in admin panel for configuration management
- Fast loading optimized with memory caching (configurable TTL)
