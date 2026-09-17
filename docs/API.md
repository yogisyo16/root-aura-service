# API Reference

Base URL (local dev, per `.env`): `http://localhost:65510`

All routes are mounted under `/api`, versioned as `/v1` or `/v2`. There is no authentication on any route yet — see [`BACKLOG.md`](./BACKLOG.md).

> **Note on response shapes:** the handlers were built incrementally and are not fully consistent yet. Some return `{"code": ..., "data": {"items": ...}}`, others return `{"Msg": ..., "Code": ...}` (capitalized, from the shared `Response` struct), and a couple (`getTodoByID`, `getTodoDetailsByID`, `getTodoDetailsByTodoId`) return the raw object with no envelope at all. This is called out per-endpoint below and tracked in the backlog.

---

## Health Check

| Method | Path | Handler |
|---|---|---|
| GET | `/api/v1/healthcheck` | `HealthCheck` |
| GET | `/api/v2/healthcheck` | `HealthCheck` |

Response `200`:
```json
{ "Msg": "Health Check", "Code": 200 }
```

---

## Users

| Method | Path | Handler |
|---|---|---|
| POST | `/api/v1/users/create` | `UserHandler.insertUser` |
| GET | `/api/v1/users` | `UserHandler.getAllUsers` |
| GET | `/api/v1/users/{id}` | `UserHandler.getUserByID` |

### `POST /users/create`
Request body:
```json
{ "first_name": "Jane", "last_name": "Doe", "email": "jane@example.com", "password": "plaintext" }
```
The password is bcrypt-hashed server-side before storage (`bcrypt.DefaultCost`) — never send a pre-hashed value.

Response `201`:
```json
{ "Msg": "Successfully Created User", "Code": 201 }
```
Errors: `400` invalid JSON body, `500` hashing failure or DB insert failure.

### `GET /users`
Response `200`:
```json
{ "code": 200, "data": { "items": [ { "id": "...", "first_name": "...", "last_name": "...", "email": "...", "password": "<bcrypt hash>", "created_at": "...", "update_at": "..." } ] } }
```
⚠️ **Known issue**: the bcrypt hash is currently returned in the `password` field. See `BACKLOG.md`.

### `GET /users/{id}`
`{id}` is the Mongo `ObjectID` hex string. Response `200`:
```json
{ "code": 200, "data": { "items": { "id": "...", "first_name": "...", "...": "..." } } }
```
`404`-style failure currently just logs and returns an empty `200` body (no explicit error status is written) — see `BACKLOG.md`.

---

## Todos

| Method | Path | Handler |
|---|---|---|
| GET | `/api/v1/todos` | `TodoHandler.getTodos` |
| GET | `/api/v1/todos/{id}` | `TodoHandler.getTodoByID` |
| POST | `/api/v1/todos/create` | `TodoHandler.createTodo` |
| PUT | `/api/v1/todos/update/{id}` | `TodoHandler.updateTodo` |
| PATCH | `/api/v1/todos/{id}/complete` | `TodoHandler.toggleComplete` |
| DELETE | `/api/v1/todos/delete/{id}` | `TodoHandler.deleteTodo` |

### `GET /todos`
Query params (all optional):
- `sort_by` — `task` \| `date_start` \| `date_due` \| `completed` \| `updated_at` \| `created_at` (default)
- `sort_order` — `ASC` (default) \| `DESC`
- `page` — default `1`
- `limit` — default `10`, max `100`

Each todo in the response is enriched with its `todo_details` (looked up by `todo_id`, `null` if none exist yet).

Response `200`:
```json
{
  "code": 200,
  "data": {
    "items": [ { "id": "...", "user_id": "...", "task": "...", "date_start": null, "date_due": null, "completed": false, "todo_details": null, "created_at": "...", "update_at": "..." } ],
    "pagination": { "current_page": 1, "total_pages": 1, "total_items": 1, "items_per_page": 10, "has_next": false, "has_prev": false }
  }
}
```

### `GET /todos/{id}`
Returns the todo (with `todo_details` embedded) **without** the `{code, data}` envelope — a raw object. `404` if not found.

### `POST /todos/create`
Request body:
```json
{ "task": "Buy milk", "date_start": "2026-09-20", "date_due": "2026-09-21T18:00:00Z", "completed": false }
```
Dates accept `YYYY-MM-DD`, full RFC3339, or `YYYY-MM-DDTHH:MM:SS`; both are optional, but if both are given `date_start` must not be after `date_due`. `task` is required.

Response `201`:
```json
{ "msg": "Successfully Created Todo", "code": 201, "data": { "id": "...", "task": "...", "...": "..." } }
```
Note: `user_id` is never set on create — see `BACKLOG.md` (todos aren't yet scoped to a user).

### `PUT /todos/update/{id}`
Same body shape as create (full replace of `task`/`date_start`/`date_due`/`completed`). Response `200`: `{ "Msg": "Successfully Updated Todo", "Code": 200 }`.

### `PATCH /todos/{id}/complete`
No body. Flips the todo's `completed` boolean. Response `200`: `{ "Msg": "Successfully toggled completion status", "Code": 200 }`.

### `DELETE /todos/delete/{id}`
Response `200` on success; on error returns HTTP `304` with `{ "Msg": "Error", "Code": 304 }` (unconventional status code choice — normally `500`/`404`).

---

## Todo Details

| Method | Path | Handler |
|---|---|---|
| GET | `/api/v1/todos/tododetails` | `TodoDetailsHandler.getTodoDetails` |
| GET | `/api/v1/todos/tododetails/{id}` | `TodoDetailsHandler.getTodoDetailsByID` |
| GET | `/api/v1/todos/tododetails/todoid/{todo_id}` | `TodoDetailsHandler.getTodoDetailsByTodoId` |
| POST | `/api/v1/todos/{todo_id}/details` | `TodoDetailsHandler.createTodoDetails` |
| DELETE | `/api/v1/todos/tododetails/delete/{id}` | `TodoDetailsHandler.deleteTodoDetails` |

### `GET /todos/tododetails`
Response `200`: `{ "code": 200, "data": { "items": [ {...TodoDetails} ] } }`.

### `GET /todos/tododetails/{id}`
Raw `TodoDetails` object (no envelope). `404` if not found.

### `GET /todos/tododetails/todoid/{todo_id}`
Raw `TodoDetails` object. If no details exist yet for that todo, returns `200` with an **empty object** (`{}`) rather than `404` — this is intentional (a todo without details yet is not an error).

### `POST /todos/{todo_id}/details`
Request body:
```json
{ "task_details": "...", "notes_details": "...", "status_details": "todo", "priority_details": "high" }
```
`{todo_id}` comes from the URL, not the body. Response `201`: `{ "Msg": "Successfully Created Todo Details", "Code": 201 }`.

### `DELETE /todos/tododetails/delete/{id}`
Response `200` on success; `304` with `{ "Msg": "Error", "Code": 304 }` on failure (same non-standard pattern as todo delete).

---

## Example requests

```bash
# Register
curl -X POST http://localhost:65510/api/v1/users/create \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Jane","last_name":"Doe","email":"jane@example.com","password":"secret123"}'

# List todos, newest first, page 2
curl "http://localhost:65510/api/v1/todos?sort_by=created_at&sort_order=DESC&page=2&limit=5"

# Create a todo
curl -X POST http://localhost:65510/api/v1/todos/create \
  -H "Content-Type: application/json" \
  -d '{"task":"Buy milk","date_due":"2026-09-20"}'

# Add details to a todo
curl -X POST http://localhost:65510/api/v1/todos/<todo_id>/details \
  -H "Content-Type: application/json" \
  -d '{"task_details":"2% milk, 1 gallon","status_details":"todo","priority_details":"low"}'
```
