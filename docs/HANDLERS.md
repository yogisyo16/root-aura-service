# Handlers

The `handlers` package is the HTTP layer — decodes requests, calls a `services` method, encodes the response. Built on [`chi`](https://github.com/go-chi/chi) for routing.

## `handlers/router.go`

Owns route wiring and cross-cutting config, not business logic.

- **`Response` struct** — `{ Msg string; Code int }`, the shared "simple" response envelope used by most write endpoints (create/update/delete). Read endpoints tend to use an ad-hoc inline `{code, data: {items}}` struct instead — see the inconsistency note in `API.md`.
- **`CreateRouter(todoHandler, userHandler, todoDetailsHandler) *chi.Mux`** — builds the router:
  - CORS via `go-chi/cors`: any origin (`https://*`, `http://*`), credentials allowed, standard methods, `Authorization` header explicitly allowed (already in place for when JWT auth lands).
  - Mounts everything under `/api`, then `/v1` (the working API) and `/v2` (currently just a health check placeholder for future versioning).
  - No middleware is attached yet — no auth, no request logging, no rate limiting.

## `handlers/userHandler.go` — `UserHandler`

Wraps `services.User`.

| Function | Route | Behavior |
|---|---|---|
| `insertUser` | `POST /users/create` | Decodes body → bcrypt-hashes the password (`bcrypt.GenerateFromPassword`, default cost) → calls `Service.InsertUser`. This is the one place password hashing happens; the service layer just stores whatever it's given. |
| `getAllUsers` | `GET /users` | Calls `Service.GetAllUsers`, wraps in `{code, data: {items}}`. |
| `getUserByID` | `GET /users/{id}` | Reads `id` via `chi.URLParam`, calls `Service.GetUserByID`. On error, logs and returns — **note**: this leaves the response unwritten with an implicit `200` (chi/net-http defaults to 200 if `WriteHeader` is never called), so a lookup failure doesn't currently surface as a proper HTTP error. Worth fixing alongside the auth work.

## `handlers/todoHandler.go` — `TodoHandler`

Wraps both `services.Todo` and `services.TodoDetails` (`DetailsService`), because todo reads are enriched with their details inline.

| Function | Route | Behavior |
|---|---|---|
| `HealthCheck` | `GET /healthcheck` (v1 & v2) | Package-level function, not a method — doesn't touch any service. |
| `getTodos` | `GET /todos` | The most involved handler: parses `sort_by`/`sort_order`/`page`/`limit` query params (with defaults and bounds-checking), fetches all todos, looks up each one's `TodoDetails` by `todo_id`, sorts in-memory (`sort.Slice`), then paginates the already-sorted slice. All sorting/pagination happens in Go after loading the *entire* collection — fine at small scale, won't scale to a large `todos` collection (no DB-level `$sort`/`$skip`/`$limit`). |
| `getTodoByID` | `GET /todos/{id}` | Fetches the todo, attaches its details if present, returns the raw object (no envelope). |
| `parseDateTime` | (helper, not a route) | Tries RFC3339, then `YYYY-MM-DD`, then `YYYY-MM-DDTHH:MM:SS`, in that order. Shared by create and update. |
| `createTodo` | `POST /todos/create` | Validates `task` is non-empty, parses optional dates, validates `date_start <= date_due` when both given, inserts. Does **not** set `UserID` on the created todo (see `BACKLOG.md`). |
| `updateTodo` | `PUT /todos/update/{id}` | Same validation as create, full-field replace. |
| `toggleComplete` | `PATCH /todos/{id}/complete` | Reads the current todo, flips `Completed`, writes it back via the same full-update method used by `updateTodo`. |
| `deleteTodo` | `DELETE /todos/delete/{id}` | Deletes; on failure responds `304` (unconventional — typically `404`/`500`). |

`TodoWithDetails` is the response-shaping struct that embeds a `*services.TodoDetails` alongside the todo fields — it's what `getTodos`/`getTodoByID` actually serialize, not the bare `services.Todo`.

## `handlers/todoDetailsHandler.go` — `TodoDetailsHandler`

Wraps `services.TodoDetails`.

| Function | Route | Behavior |
|---|---|---|
| `getTodoDetails` | `GET /todos/tododetails` | All details, `{code, data: {items}}` envelope. |
| `getTodoDetailsByID` | `GET /todos/tododetails/{id}` | Raw object, `404` if not found. |
| `getTodoDetailsByTodoId` | `GET /todos/tododetails/todoid/{todo_id}` | Raw object; returns `200` with an empty object if the todo has no details yet (matches the service-layer "not found is not an error" behavior). |
| `createTodoDetails` | `POST /todos/{todo_id}/details` | `todo_id` comes from the URL, not the body; decodes the rest from `CreateTodoDetailsRequest`. |
| `deleteTodoDetails` | `DELETE /todos/tododetails/delete/{id}` | Same `304`-on-error pattern as `deleteTodo`. |

## Request/response shaping structs

Handlers define their own request DTOs rather than reusing the `services` structs directly for input (`CreateTodoRequest`, `UpdateTodoRequest`, `CreateTodoDetailsRequest`), which keeps the wire format decoupled from the Mongo document shape. Output, though, is inconsistent — some handlers reuse `services.X` directly, `todoHandler.go` defines its own `TodoWithDetails` for enriched responses, and there's no single shared "list response" or "error response" builder, which is why the envelope shapes vary across endpoints (documented per-endpoint in `API.md`).
