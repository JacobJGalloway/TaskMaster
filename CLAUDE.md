# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

Task Master — a small full-stack task-tracking app, built as a demonstration project. Its purpose isn't to be a polished product; it's to show working, hands-on knowledge of **current Go** and **current Vue.js**, and how the two sides talk to each other to form one application.

**Services:**
1. `TaskMaster-api/` — Go backend. REST API backed by SQLite, built with Go's standard library only (`net/http`, `database/sql`) — no third-party web framework or ORM.
2. `TaskMaster-web/` — Vue 3 frontend. Vite + TypeScript SPA that calls the API over plain `fetch` — no Axios, no generated client.

**Database:** SQLite, via the pure-Go `modernc.org/sqlite` driver — no separate SQLite server/install required. The `db.sqlite` file lives in `TaskMaster-api/` and is created automatically on first run. Schema (`CREATE TABLE IF NOT EXISTS ...`) lives directly in `main.go` and is applied on every startup; there's no separate migration tool, which is fine at this scale but the first thing that would need to change if this grew past a demo.

## Build & Run

```bash
cd TaskMaster-api
go run main.go          # API on :8080, creates/opens db.sqlite, applies schema

cd TaskMaster-web
npm install              # first run only
npm run dev               # Vite dev server on :5173
npm run build               # vue-tsc type-check + production build
npm run test:unit             # Vitest unit tests
npm run lint                    # oxlint + eslint, auto-fix
npm run format                    # Prettier
```

Run both in separate terminals and leave both running. The frontend hardcodes `http://localhost:8080` as the API base URL (no `.env` yet).

## Architecture

### TaskMaster-api layers

```
main.go
  └─ handler    (HTTP: parse request, call service, write response)
       └─ service   (business logic / orchestration)
            └─ repository  (SQL: talk to the database)
                 └─ model      (plain data structs shared by every layer)
```

- `model/` — plain structs with JSON tags (`User`, `StudyTask`). No behavior, just shape.
- `repository/` — one file per entity (`task_repository.go`, `user_repository.go`). Small interface (e.g. `TaskRepository`) + concrete `sqlite*Repository` implementing it with raw, parameterized SQL via `database/sql`. The *only* layer that knows SQL exists.
- `service/` — one file per entity. Thin pass-through today (`task_service.go` mostly delegates to the repository) — this is the seam where business rules/validation would go as the app grows.
- `handler/` — one file per entity. Decodes the JSON request body, calls the matching service method, encodes the JSON response or writes an HTTP error. The *only* layer that knows HTTP exists.

Each layer depends on the layer below it through an **interface** (`service.TaskService`, `repository.TaskRepository`), never a concrete type. `main.go` constructs the concrete implementations and injects them downward — plain constructor-based dependency injection, no framework/container.

### Routes

Registered on a single `http.ServeMux` in `main.go`, using Go 1.22+'s method-aware routing patterns — no router library needed:

| Method & path | Purpose |
|---|---|
| `GET /health` | Liveness check |
| `POST /users` / `GET /users` / `GET /users/{id}` | Create / list / fetch a user |
| `GET /users/{id}/tasks` | List a user's **incomplete** tasks only |
| `POST /tasks` | Create a task |
| `POST /tasks/{id}/complete` | Mark a task complete (sets `completed = 1`, stamps `completed_at`) |

CORS is handled by a small `withCORS` middleware wrapping the mux, allowing only `http://localhost:5173` (the Vite dev origin) and short-circuiting `OPTIONS` preflight.

### TaskMaster-web structure

```
src/
  App.vue                 - root; owns selectedUserId, passes it down to TaskList
  components/
    UserPicker.vue          - list/select/create users
    TaskList.vue              - list/create/complete tasks for the selected user
  composables/
    useUsers.ts                - fetch/create users; owns its own users/error/isLoading refs
    useTasks.ts                   - fetch/create/complete tasks; owns its own reactive state
  types/
    User.ts, StudyTask.ts            - TypeScript interfaces mirroring the Go JSON structs
  stores/counter.ts, router/index.ts    - Vue project-init scaffold, not actively used
```

**Composables over a central store.** Pinia is installed but state here is managed with two small composables (`useUsers`, `useTasks`) instead — plain functions returning reactive `ref`s. Each call site gets its own independent state; nothing is shared globally. Deliberate scope decision for an app this size.

**Mutation pattern.** Writes (`createTask`, `completeTask`, `createUser`) re-fetch the list afterward (`await fetchTasks(userId)`) rather than optimistically patching local state — the UI's source of truth after any mutation is "ask the server again," not "assume success and update the array locally."

### Known gaps (intentional, not oversight)

These are deliberate scope cuts for a demo, each with a documented reason — see the Backlog section in the root README for the current prioritized list of what's next:

- **No auth** — no login/session layer; any client can call any endpoint.
- **No ORM** — all SQL is hand-written in the repository layer (control/clarity at this scale vs. boilerplate as it grows).
- **Coarse error handling** — handlers mostly return generic `400`/`500`s rather than typed, specific error codes.
- **Hardcoded API base URL** — frontend hardcodes `http://localhost:8080`; no `.env`-driven config yet.
