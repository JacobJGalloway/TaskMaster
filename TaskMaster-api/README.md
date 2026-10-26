# TaskMaster-api — technical overview

This is the Go backend for Task Master. It's a REST API, backed by SQLite, built with Go's standard library only (`net/http`, `database/sql`) — no third-party web framework or ORM. For the non-technical project summary and run instructions, see the [root README](../README.md).

## Layered architecture

The code is split into four layers, each with one job, wired together in `main.go`:

```
main.go
  └─ handler   (HTTP: parse request, call service, write response)
       └─ service   (business logic / orchestration)
            └─ repository   (SQL: talk to the database)
                 └─ model   (plain data structs shared by every layer)
```

- **`model/`** — plain structs with JSON tags (`User`, `StudyTask`). No behavior, just shape. These are what get marshaled to/from JSON on the wire and scanned to/from SQLite rows.
- **`repository/`** — one file per entity (`task_repository.go`, `user_repository.go`). Each defines a small interface (e.g. `TaskRepository`) and a concrete `sqlite*Repository` struct implementing it with raw SQL via `database/sql`. This is the *only* layer that knows SQL exists.
- **`service/`** — one file per entity. Thin for now (`task_service.go` mostly just delegates to the repository), but this is the seam where business rules would go if the app grew — e.g. validation, authorization, combining multiple repository calls.
- **`handler/`** — one file per entity. Decodes the JSON request body (if any), calls the matching service method, and encodes the JSON response (or writes an HTTP error). This is the *only* layer that knows HTTP exists.

Each layer only depends on the one below it, and depends on it through an **interface**, not a concrete type (`service.TaskService`, `repository.TaskRepository`). `main.go` is what actually constructs the concrete implementations and injects them downward — plain constructor-based dependency injection, no framework/container involved.

Why it's built this way: it keeps HTTP concerns (status codes, JSON) out of the SQL code, and keeps SQL out of the HTTP code. In a bigger app it also makes each layer independently testable — you can fake a `TaskRepository` to unit test `TaskService`, or fake a `TaskService` to unit test `TaskHandler`, without touching a real database.

## Entities

| Entity | Fields | Notes |
|---|---|---|
| `User` | `id`, `name` | Minimal — just enough to own a set of tasks |
| `StudyTask` | `id`, `user_id`, `title`, `duration_minutes`, `period_start`, `period_end`, `completed`, `completed_at` | `completed_at` is a nullable pointer (`*string`), set only when a task is marked complete |

The schema (`CREATE TABLE IF NOT EXISTS ...`) lives directly in `main.go` and is applied against `db.sqlite` on every startup — there's no separate migration tool, which is fine at this scale but is the first thing that would need to change if this grew past a demo.

## Routes

Registered on a single `http.ServeMux` in `main.go`, using Go 1.22+'s method-aware routing patterns (`"GET /users"`, `"POST /tasks/{id}/complete"`, etc. — no router library needed):

| Method & path | Handler | Purpose |
|---|---|---|
| `GET /health` | inline in `main.go` | Liveness check, returns `OK` |
| `POST /users` | `UserHandler.Create` | Create a user |
| `GET /users` | `UserHandler.GetUsers` | List all users |
| `GET /users/{id}` | `UserHandler.GetUserByID` | Fetch one user |
| `GET /users/{id}/tasks` | `TaskHandler.ListActiveByUser` | List a user's **incomplete** tasks |
| `POST /tasks` | `TaskHandler.Create` | Create a task |
| `POST /tasks/{id}/complete` | `TaskHandler.MarkComplete` | Mark a task complete (sets `completed = 1`, stamps `completed_at`) |

Note `ListActiveByUser` only returns tasks where `completed = 0` — completed tasks fall out of the list rather than being deleted (the row still exists in SQLite, just filtered out of this particular query).

CORS is handled by a small `withCORS` middleware wrapping the whole mux, allowing requests from `http://localhost:5173` (the Vite dev server origin) and short-circuiting `OPTIONS` preflight requests.

## Data flow, end to end

Using "create a task" as the representative example:

1. Frontend sends `POST /tasks` with a JSON body.
2. `main.go`'s mux routes it to `TaskHandler.Create`.
3. The handler decodes the JSON body into a `model.StudyTask` struct — this is the HTTP/JSON boundary.
4. The handler calls `taskService.CreateTask(ctx, &task)` — crossing from the HTTP layer into business logic.
5. The service (today, a pass-through) calls `taskRepo.Create(ctx, t)`.
6. The repository runs a parameterized `INSERT` against `db.sqlite` via `database/sql`, using `?` placeholders (never string-concatenated SQL, to avoid injection).
7. Success/failure bubbles back up through service → handler, which writes `201 Created` or a `5xx` with an error message.

Reads (e.g. `GET /users/{id}/tasks`) run the same layers in reverse: repository runs a `SELECT`, scans rows into `[]model.StudyTask`, and the handler `json.Encode`s that slice straight into the HTTP response body — so the database row shape and the JSON wire shape are the same struct, by design, for simplicity.

## Local dependencies worth knowing for Round 2

- **`modernc.org/sqlite`** — a pure-Go, CGo-free SQLite driver. This is why there's no "install SQLite" step for local dev: Go itself can read/write the `.sqlite` file.
- **No ORM** — all SQL is hand-written in the repository layer. Worth being able to explain the tradeoff (more control/clarity at this scale vs. more boilerplate as it grows).
- **No auth** — there's no login/session layer; any client can call any endpoint. Worth having a view on how you'd add it (middleware around the mux, a new layer, JWT vs. sessions, etc.) as likely follow-up territory.
- **Error handling is coarse** — handlers mostly return generic `400`/`500`s rather than typed, specific error codes. A reasonable thing to be asked to extend.

## Running it

```bash
cd TaskMaster-api
go run main.go
```

Starts on `:8080`, creates/opens `db.sqlite` in this folder, and applies the schema on boot. See the [root README](../README.md) for running this alongside the frontend.
