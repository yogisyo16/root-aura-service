# Services

The `services` package is the data-access layer — one file per domain, each owning a Mongo collection, a struct (which doubles as the Mongo document shape and the JSON wire shape), and the CRUD methods handlers call into. There's no repository interface separating Mongo specifics from business logic; handlers call these structs' methods directly.

## `services/main.go` — `Models`

```go
type Models struct {
    Todo    Todo
    User    User
    Details TodoDetails
}
```
A simple aggregate struct. In practice `api/main.go` constructs `Todo`, `User`, and `TodoDetails` separately and passes them straight to the handler constructors — `Models` itself isn't currently used to wire anything, it just documents "these are the domain models."

## `services/todoServices.go` — `Todo`

Collection: `todos_db.todos`.

```go
type Todo struct {
    ID        string     `bson:"_id,omitempty"`
    UserID    string     `bson:"_user_id"`
    Task      string     `bson:"_task"`
    DateStart *time.Time `bson:"_date_start,omitempty"`
    DateDue   *time.Time `bson:"_date_due,omitempty"`
    Completed bool       `bson:"_completed"`
    CreatedAt time.Time  `bson:"_created_at,omitempty"`
    UpdatedAt time.Time  `bson:"_update_at,omitempty"`
}
```
`UserID` exists on the model but nothing currently sets or filters by it — see `BACKLOG.md`.

| Method | Purpose |
|---|---|
| `New(mongo *mongo.Client) Todo` | Stores the shared Mongo client in the package-level `client` var and returns a zero-value `Todo` service. Called once at startup. |
| `GetAllTodos() ([]Todo, error)` | Returns every document in `todos`, unfiltered. |
| `GetTodoById(id string) (Todo, error)` | Looks up by hex `ObjectID`. |
| `InsertTodo(entry Todo) (*Todo, error)` | Inserts, stamping `CreatedAt`/`UpdatedAt`; returns the inserted doc with its generated `ID` populated. |
| `UpdatedTodo(id string, entry Todo) (*mongo.UpdateResult, error)` | Full-field `$set` update (task, dates, completed, `UpdatedAt`). |
| `DeleteTodo(id string) error` | Deletes by `ObjectID`. |

Note: the package-level `client *mongo.Client` and `returnTodosCollection` helper are declared here and reused (via shared package state) by `userServices.go` and `todoDetailsServices.go` too — all three services share one Mongo client set once via `services.New()`.

## `services/userServices.go` — `User`

Collection: `todos_db.users`.

```go
type User struct {
    ID        string    `bson:"_id,omitempty"`
    FirstName string    `bson:"_first_name,omitempty"`
    LastName  string    `bson:"_last_name,omitempty"`
    Email     string    `bson:"_email,omitempty"`
    Password  string    `bson:"_password,omitempty"` // bcrypt hash, see note below
    CreatedAt time.Time `bson:"_created_at,omitempty"`
    UpdatedAt time.Time `bson:"_update_at,omitempty"`
}
```

| Method | Purpose |
|---|---|
| `GetAllUsers() ([]User, error)` | Returns every user document, unfiltered. |
| `InsertUser(entry User) error` | Inserts a user with timestamps. Expects `Password` to already be a bcrypt hash — hashing happens in the handler, not here. |
| `GetUserByID(id string) (User, error)` | Looks up by hex `ObjectID`. |

⚠️ `Password` is tagged `json:"password,omitempty"`, so it round-trips into API responses as the bcrypt hash. Should be `json:"-"`. Tracked in `BACKLOG.md`.

There's no `GetUserByEmail` yet — required for a login flow, also tracked in the backlog.

## `services/todoDetailsServices.go` — `TodoDetails`

Collection: `todos_db.todo_details`. A one-to-one(ish) extension of `Todo`, joined by `todo_id`.

```go
type TodoDetails struct {
    ID              string    `bson:"_id,omitempty"`
    TodoID          string    `bson:"_todo_id,omitempty"`
    TaskDetails     string    `bson:"_task_details"`
    NotesDetails    string    `bson:"_notes_details"`
    StatusDetails   string    `bson:"_status_details"`
    PriorityDetails string    `bson:"_priority_details"`
    CreatedAt       time.Time `bson:"_created_at,omitempty"`
    UpdatedAt       time.Time `bson:"_updated_at,omitempty"`
}
```

| Method | Purpose |
|---|---|
| `NewTodoDetailsService(mongo *mongo.Client) TodoDetails` | Same pattern as `todoServices.New` — stores the client, returns a zero-value service. |
| `GetAllTodosDetails() ([]TodoDetails, error)` | All documents, unfiltered. |
| `GetTodoDetailsById(id string) (TodoDetails, error)` | Lookup by the details doc's own `ObjectID`. |
| `GetTodoDetailsByTodoId(todoId string) (TodoDetails, error)` | Lookup by the **parent todo's** id. Returns a zero-value `TodoDetails{}` with a `nil` error when nothing is found (not treated as an error) — callers check `details.ID != ""` to know whether details exist. |
| `InsertTodoDetails(entry TodoDetails) error` | Inserts with timestamps. |
| `UpdateTodoDetails(id string, entry TodoDetails) (*mongo.UpdateResult, error)` | Full-field `$set` update. |
| `DeleteTodoDetails(id string) error` | Deletes by `ObjectID`. |

## `services/authServices.go` — not yet implemented

Currently an empty file (just `package services`). Planned to hold `GenerateToken`/`ValidateToken` for JWT auth — see `BACKLOG.md` for the full plan.
