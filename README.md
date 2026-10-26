# Task Master

Task Master is a small full-stack task-tracking app, built as a demonstration project. Its purpose isn't to be a polished product — it's to show working, hands-on knowledge of **current Go** and **current Vue.js**, and how the two sides talk to each other to form one application.

A user picks (or creates) themselves from a list, then adds "tasks" — a title, a duration, and a start/end time window — and checks them off as they're completed. Simple on purpose: the goal is to demonstrate clean structure and correct integration, not feature breadth.

This README is the plain-language overview. For a deeper, component-by-component technical breakdown of each side, see:

- [`TaskMaster-api/README.md`](./TaskMaster-api/README.md) — backend internals
- [`TaskMaster-web/README.md`](./TaskMaster-web/README.md) — frontend internals

## How it's organized

The project is two independent applications that live side by side and talk to each other over HTTP:

| Folder | What it is |
|---|---|
| `TaskMaster-api` | The **backend** — a Go program that stores data and exposes it over a REST API |
| `TaskMaster-web` | The **frontend** — a Vue.js single-page app that runs in the browser and calls that API |

They're started and run separately (see below), and during local development they talk to each other across two different "ports" on the same machine — `localhost:8080` for the API, `localhost:5173` for the web app.

## Tech stack

### Backend — `TaskMaster-api`

| Piece | Role |
|---|---|
| **Go** | The language the backend is written in |
| **Go's standard `net/http` package** | Handles incoming web requests — no third-party web framework was used, to demonstrate familiarity with the standard library itself |
| **SQLite** (via the `modernc.org/sqlite` driver) | The database — a single file on disk, no separate database server to install or run |
| **`database/sql`** | Go's standard interface for talking to the database |

Supporting libraries pulled in automatically alongside the SQLite driver (these aren't used directly, they're dependencies of the database driver): `github.com/google/uuid`, `github.com/dustin/go-humanize`, `github.com/mattn/go-isatty`, `github.com/ncruces/go-strftime`, and a few `modernc.org/*` helper packages.

### Frontend — `TaskMaster-web`

| Piece | Role |
|---|---|
| **Vue 3** | The JavaScript framework the UI is built with |
| **Vite** | The dev server and build tool — gives fast reloads while developing |
| **TypeScript** | Adds type-checking on top of JavaScript |
| **Pinia** | Vue's official state-management library (included for completeness; most state here is handled with simpler Vue composables) |
| **Vue Router** | Client-side page routing |
| **Font Awesome** (`@fortawesome/vue-fontawesome` + icon packs) | Icons used in the UI, such as the "mark complete" checkmark |
| **Vitest** | Unit test runner |
| **ESLint / oxlint / Prettier** | Code linting and formatting |

## Prerequisites

Install these before running the project locally:

1. **[Go](https://go.dev/dl/)** — version 1.26 or newer (the backend is built against this version)
2. **[Node.js](https://nodejs.org/)** — version 22.18 or newer (or Node 24.12+); this also provides `npm`
3. **SQLite** — technically, you don't need to install a separate SQLite server or CLI. The Go driver used here (`modernc.org/sqlite`) is a pure-Go implementation that reads/writes the `db.sqlite` file directly, with no external SQLite installation required. If you'd like to inspect the database file by hand, installing the [SQLite command-line tool](https://www.sqlite.org/download.html) is optional but handy.
4. **Git** — to clone/manage the repository
5. *(Optional but recommended)* **[VS Code](https://code.visualstudio.com/)** with the Go extension and the Vue (Volar) extension, since both READMEs assume you may be running things from VS Code's integrated terminal

No other installs are required — the backend has no external database server, message queue, or cache to stand up.

## Running the application locally

The backend and frontend are two separate, long-running processes. Start each **in its own terminal session** (two PowerShell/Bash windows, or two VS Code integrated terminal tabs) and leave both running while you work.

### 1. Start the backend API

In the first terminal:

```bash
cd TaskMaster-api
go run main.go
```

You should see log output confirming the database schema was initialized, followed by the server listening. The API is now available at `http://localhost:8080` (try `http://localhost:8080/health` in a browser — it should return `OK`).

On first run, Go will create the `db.sqlite` file in that folder automatically — nothing to set up manually.

### 2. Start the frontend web app

In a **second** terminal, leaving the first one running:

```bash
cd TaskMaster-web
npm install
npm run dev
```

`npm install` only needs to be run once (or whenever dependencies change). After that, `npm run dev` starts Vite's dev server, and it will print a local URL — typically `http://localhost:5173`. Open that in your browser.

### Doing this from VS Code

If you'd rather use VS Code instead of standalone terminal windows:

1. Open the `TaskMaster` folder in VS Code.
2. Open a new integrated terminal (`` Ctrl+` ``), `cd TaskMaster-api`, and run `go run main.go`.
3. Click the **"+"** in the terminal panel to open a second, independent terminal tab, `cd TaskMaster-web`, and run `npm run dev`.
4. Both terminals keep running side by side as separate tabs — closing one does not stop the other.

Once both are running, the web app in your browser will make requests to the API in the background, and you should be able to add a user, add tasks, and mark them complete.
