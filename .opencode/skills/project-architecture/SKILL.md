---
name: project-architecture
description: Use when planning, implementing, reviewing, or documenting changes in gitlab-board to preserve backend/frontend boundaries, data flows, wiring, contracts, persistence, deployment, and required verification.
compatibility: opencode
metadata:
  project: gitlab-board
  audience: development-agents
---

# Project Architecture

Use this skill as an architecture navigation and impact-analysis procedure. The
repository documents and current source remain the source of truth; do not rely
on this summary when they disagree.

## Required Reading

Always read `AGENTS.md` first. Then select documents by impact:

- Read `.docs/architecture.md` for cross-application flows, runtime topology,
  storage, GitLab integration, deployment, or configuration.
- Read `.docs/arch-backend.md` and `.docs/rules.md` for any backend change.
- Read `.docs/arch-frontend.md` and `.docs/rules.md` for any frontend change.
- Read all four architecture documents for API DTO, auth, WebSocket, shared
  data flow, Docker/env, or changes spanning frontend and backend.

After reading documentation, inspect the actual entrypoints, wiring, callers,
implementations, and tests touched by the request. Source code and the toolchain
files named in `AGENTS.md` outrank stale prose.

## System Map

- `backend/` is one Go/Gin process containing REST, WebSocket, and the periodic
  GitLab sync worker. Application operations sit between Gin controllers and
  domain/persistence/external ports. GORM supports SQLite, PostgreSQL, and
  MySQL.
- `frontend/board/` is an Angular SPA. State is distributed across components,
  URL query parameters, RxJS services, REST services, and a sprint WebSocket
  client; there is no global store.
- GitLab GraphQL handles bulk synchronization; REST/OAuth handles user
  operations. Several writes are local-first with best-effort GitLab
  write-back, so partial failure is an explicit design concern.
- WebSocket rooms and online state are process-local. There is no durable queue
  or shared realtime state across backend replicas.

## Boundary Checks

For each proposed or actual change, trace the complete relevant path:

```text
Frontend route/component -> service -> HTTP/WebSocket contract
-> Gin controller/middleware -> use case/service -> domain/port
-> GORM or GitLab adapter -> database/external system
```

Only include layers that participate in the behavior.

Backend checks:

- Controllers own HTTP parsing, DTOs, status codes, and middleware wiring.
- Use cases/services own orchestration and business operations; domain models
  own invariants; repositories/specifications own persistence/query details.
- New external details should be behind a narrow app/domain port when that
  prevents business logic from depending on infrastructure.
- A controller requires both DI registration and router registration. A new
  persistent model requires repository wiring and `ALL_MODELS` registration.
- Mutation routes require explicit `AuthRequired` or `AdminRequired`; global
  JWT parsing alone does not authorize access.
- Binding changes must keep persistence, history, and Sprint WebSocket events
  consistent.
- SQL/search/report behavior must account for all supported database dialects.
  Multi-step atomic behavior requires a real transaction boundary.

Frontend checks:

- Keep route shells, kanban features, reusable UI, and singleton core
  infrastructure in their documented modules.
- Obtain backend URLs from `CoreConfigService`; client role visibility is not
  an authorization boundary.
- Respect current NgModule versus standalone component conventions and use
  public Angular package entrypoints.
- Define cache invalidation, subscription lifetime, optimistic rollback, query
  parameter persistence, and WebSocket reconciliation where relevant.
- Contract changes must update both backend and frontend in one vertical slice.

## Architecture Impact Classification

Treat a change as architecture-relevant when it modifies any of these:

- module ownership, layer boundaries, or dependency direction;
- request/event payloads, auth rules, or error/status semantics;
- runtime processes, background work, realtime topology, or failure behavior;
- persistence models, migrations/`AutoMigrate`, transactions, queries, or
  supported database behavior;
- DI/router/module wiring, external integrations, configuration, env loading,
  Docker/nginx behavior, or build/runtime topology;
- mandatory development, test, build, or operational commands.

Local implementation details that preserve these contracts usually do not
require architecture documentation changes. Do not turn architecture documents
into a file-by-file changelog.

## Role-Specific Use

- Planner: identify affected boundaries and documentation before decomposing
  vertical implementation units.
- Implementer: reread the relevant rules before editing and report every
  architecture impact signal; leave architecture document edits to the final
  documentation agent.
- Reviewer: verify requirements and boundary rules against actual code, then
  explicitly hand documentation signals to the final phase.
- Documentation agent: compare accepted final code with all signaled areas and
  update only documents whose architectural facts changed. An explicit
  `NO_ARCH_CHANGE` is a valid result after inspection.

## Verification Routing

- Backend-only: run targeted integration coverage, then the required SQLite
  integration suite. Add PostgreSQL/MySQL suites for SQL, search, report, or
  dialect-sensitive changes.
- Frontend-only: run the production build and relevant unit command when
  feasible. Report the documented suite-wide Angular test compilation failures
  honestly; do not claim nonexistent lint or e2e commands.
- API DTO, auth, or WebSocket contract: verify both applications as one slice.
- Docker/env: verify development and container paths separately.

Use exact current commands from `AGENTS.md` and the relevant architecture
document rather than reconstructing them from memory.
