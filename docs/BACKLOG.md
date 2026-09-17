# Backlog

Things identified while documenting the project that aren't implemented yet, or are implemented but worth revisiting. Not in priority order beyond the numbering within each section.

## Planned: JWT authentication & authorization

Decided approach: JWT, chosen deliberately for learning value. `services/authServices.go` is currently an empty stub — nothing below is implemented yet.

1. **Add a JWT library** — `github.com/golang-jwt/jwt/v5` (actively maintained; avoid the abandoned/CVE'd `dgrijalva/jwt-go`).
2. **Add a signing secret via env var** — `JWT_SECRET` in `.env`/`Makefile`, read via `os.Getenv`. Fail fast at startup if unset.
3. **Fill in `services/authServices.go`**:
   - `GenerateToken(userID string) (string, error)` — HS256, claims limited to `sub`, `iat`, `exp` (short-lived, e.g. 30–60 min for local learning).
   - `ValidateToken(tokenString string) (userID string, err error)` — verifies signature and expiry.
4. **Add `GetUserByEmail` to `userServices.go`** — needed for login lookup; doesn't exist today (only `GetUserByID`).
5. **Add a login handler** — `POST /api/v1/auth/login`: look up by email → `bcrypt.CompareHashAndPassword` → `GenerateToken` on success. Use one generic "invalid credentials" message for both "no such user" and "wrong password" so failures don't leak which one was wrong.
6. **Add JWT middleware** — reads `Authorization: Bearer <token>`, validates it, puts the user ID on `r.Context()`, `401`s otherwise. (`router.go`'s CORS config already allows the `Authorization` header.)
7. **Wire the middleware into `router.go`** via `router.Group`/`router.With(...)`: keep `/healthcheck`, `/users/create`, `/auth/login` public; protect everything else.
8. **Scope todo queries to the authenticated user** — `Todo.UserID` already exists on the model but nothing sets or filters by it today (see below); once auth exists, `GetAllTodos`/`GetTodoById`/update/delete need to filter by the requesting user's ID, otherwise a valid token from any user can read/edit everyone's todos. This is the *authorization* half, not just authentication.
9. **Tests** — table-driven tests for `GenerateToken`/`ValidateToken` (valid, expired, tampered, wrong signature) and for the login handler (correct creds, wrong password, unknown email).

Explicitly deferred beyond this first pass: refresh tokens/rotation, role-based permissions, login rate-limiting.

## Known issues found while documenting

1. **`User.Password` leaks the bcrypt hash** — tagged `json:"password,omitempty"` in `services/userServices.go`, so `GET /users` and `GET /users/{id}` currently serialize the hash back to clients. Should be `json:"-"`.
2. **Todos aren't scoped to a user** — `Todo.UserID` exists on the model but `createTodo` never sets it and no read/update/delete filters by it. Meaningful once auth lands (see item 8 above); until then every todo is effectively global.
3. **`getUserByID` swallows lookup failures** — on error it logs and returns without calling `WriteHeader`, so a failed lookup silently responds `200` with an empty body instead of `404`.
4. **Inconsistent response envelopes** — list endpoints mostly use `{code, data: {items}}`, write endpoints mostly use the shared `Response{Msg, Code}` struct, but `getTodoByID`, `getTodoDetailsByID`, and `getTodoDetailsByTodoId` return the raw object with no envelope at all. Worth picking one convention. Full breakdown in `API.md`.
5. **Delete failures return HTTP `304`** — both `deleteTodo` and `deleteTodoDetails` respond `304 Not Modified` on error, which isn't the conventional meaning of that status (usually a caching signal, not an error). `404` (not found) or `500` (DB error) would be more standard depending on the actual failure.
6. **`getTodos` sorts and paginates in application memory** — it loads the *entire* `todos` collection, joins details per-item, sorts with `sort.Slice`, then slices for pagination. Fine at current scale; won't scale well once the collection is large (no DB-level `$sort`/`$skip`/`$limit`).
