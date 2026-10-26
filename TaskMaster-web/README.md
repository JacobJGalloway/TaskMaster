# TaskMaster-web — technical overview

This is the Vue 3 frontend for Task Master. It talks to the Go backend in `../TaskMaster-api` over plain `fetch` calls — no Axios, no generated API client, no server-side rendering. For the non-technical project summary and run instructions, see the [root README](../README.md).

## Major components

```
src/
  App.vue                     - root component; owns selectedUserId state
  components/
    UserPicker.vue             - list/select/create users
    TaskList.vue                - list/create/complete tasks for the selected user
  composables/
    useUsers.ts                 - fetch/create users; owns users/error/isLoading refs
    useTasks.ts                  - fetch/create/complete tasks; owns tasks/error/isLoading refs
  types/
    User.ts, StudyTask.ts         - TypeScript interfaces mirroring the Go JSON structs
  stores/counter.ts              - Pinia scaffold (unused placeholder from project init)
  router/index.ts                - Vue Router scaffold (single app view, not actively used for navigation)
  main.ts                        - app bootstrap: creates the Vue app, installs Pinia + router, mounts to #app
```

**Composables over a central store.** Despite Pinia being installed, state here is managed with two small "composables" (`useUsers`, `useTasks`) — plain functions that return reactive `ref`s and the functions that mutate them. Each call site that invokes `useUsers()`/`useTasks()` gets its *own* independent reactive state; nothing is shared globally. That's a deliberate scope decision for an app this size (one user list, one task list, no cross-component sharing need) rather than an oversight — worth being able to explain/contrast with "when would I reach for Pinia instead."

**Component responsibilities:**
- `UserPicker.vue` — on mount, calls `fetchUsers()` and renders the result in a `<select>`. Also has an inline "add user" form. Communicates the selected user ID *up* to `App.vue` via `defineModel` (Vue's two-way-binding macro — equivalent to emitting `update:selectedUserId`).
- `TaskList.vue` — receives `userId` as a prop from `App.vue`. On mount *and* whenever `userId` changes (via a `watch`), it calls `fetchTasks(userId)`. Renders the active task list, a "complete" button (Font Awesome checkmark icon) per task, and an inline "add task" form.
- `App.vue` — holds `selectedUserId` and passes it down to `TaskList`; otherwise has no logic of its own. This is the only place the two components are connected.

## Data flow: API response → UI state → template

This is the core pattern repeated across both composables, using `useTasks.ts` / `TaskList.vue` as the concrete example:

1. **Trigger.** `TaskList.vue`'s `onMounted` (and its `watch` on `props.userId`) calls `fetchTasks(userId)`.
2. **Request.** Inside `useTasks.ts`, `fetchTasks` sets `isLoading.value = true`, then calls `fetch('http://localhost:8080/users/{id}/tasks')` — a direct call to the Go API's `GET /users/{id}/tasks` route.
3. **Backend round-trip.** The Go API reads SQLite, returns a JSON array of task objects (see the [API README](../TaskMaster-api/README.md) for what happens on that side).
4. **Hydration.** `tasks.value = await response.json()` — the raw JSON array is assigned directly into a Vue `ref<StudyTask[]>`. This is the hydration step: the instant that assignment happens, Vue's reactivity system marks `tasks` as changed.
5. **Re-render.** Because `TaskList.vue`'s `<template>` does `v-for="task in tasks"`, and `tasks` is the same `ref` returned by `useTasks()`, Vue automatically re-renders the `<ul>` with the new data — no manual DOM manipulation, no explicit "setState" call.
6. **Error/loading state** follow the same shape: `error` and `isLoading` are refs too, set in `catch`/`finally` blocks, available for the template to branch on (e.g. to show a spinner or message) — though the current templates don't yet visibly use `error`/`isLoading`, which is a reasonable thing to be asked to add.

Writes (`createTask`, `completeTask`, `createUser`) follow the same shape in reverse, then **re-fetch** the list (`await fetchTasks(userId)` at the end of `createTask`/`completeTask`) rather than optimistically patching local state — i.e. the UI's source of truth after any mutation is always "ask the server again," not "assume the mutation succeeded and update the array locally." Worth having a view on that tradeoff (simplicity/correctness vs. an extra round trip) for Round 2.

## Request specifics worth knowing

- All API calls hardcode `http://localhost:8080` — there's no `.env`-driven base URL yet. Fine for local dev, the first thing you'd parameterize for any other environment.
- `StudyTask`'s `period_start`/`period_end` are sent as raw strings from `<input type="datetime-local">` — no date library (e.g. date-fns/Luxon) is in play; formatting/parsing is implicit.
- CORS: the API explicitly allows the Vite dev origin (`localhost:5173`); if you ever change the dev port, the API's `withCORS` middleware needs updating too (see API README).

## Tooling

- **Vite** — dev server + build tool. `npm run dev` serves with HMR; `npm run build` type-checks (via `vue-tsc`) then builds.
- **TypeScript** — `strict`-style project config split across `tsconfig.app.json` / `tsconfig.node.json` / `tsconfig.vitest.json`.
- **Vitest** + **@vue/test-utils** — unit tests live in `src/__tests__/`.
- **ESLint (+ oxlint) / Prettier** — `npm run lint` runs both linters; `npm run format` runs Prettier.

### Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur). TypeScript can't natively type-check `.vue` files, so `vue-tsc` (via Volar in the editor) stands in for plain `tsc`.

### Common commands

```bash
npm install        # install dependencies (first time only)
npm run dev         # start the dev server (http://localhost:5173)
npm run build        # type-check + production build
npm run test:unit     # run Vitest unit tests
npm run lint            # oxlint + eslint, auto-fix
npm run format           # Prettier
```

See the [root README](../README.md) for running this alongside the backend API in a second terminal/VS Code tab.
