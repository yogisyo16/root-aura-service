package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// Response is the shared "simple" reply shape used by most write endpoints
// (create/update/delete): {"Msg": "...", "Code": 200}. Read endpoints mostly
// use an ad-hoc {code, data: {items}} shape instead — see docs/API.md for
// the per-endpoint breakdown of this inconsistency.
type Response struct {
	Msg  string
	Code int
}

// CreateRouter builds the full chi router: CORS config, then every route
// under /api/v1 (the working API) and /api/v2 (currently just a health
// check placeholder for future versioning). No auth middleware is attached
// yet — every route below is public. See docs/BACKLOG.md for the planned
// JWT middleware and which routes should stay public once it lands.
func CreateRouter(todoHandler *TodoHandler, userHandler *UserHandler, todoTodoDetailsHandler *TodoDetailsHandler) *chi.Mux {
	router := chi.NewRouter()

	// Authorization is already allowed here so the CORS config won't need
	// touching once JWT auth (Authorization: Bearer <token>) is added.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CRSF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(router chi.Router) {
		router.Route("/v1", func(router chi.Router) {
			// Health Check
			router.Get("/healthcheck", HealthCheck)
			// User Routes
			router.Post("/users/create", userHandler.insertUser)
			router.Get("/users", userHandler.getAllUsers)
			router.Get("/users/{id}", userHandler.getUserByID)

			// Todo Routes
			router.Get("/todos", todoHandler.getTodos)
			router.Get("/todos/{id}", todoHandler.getTodoByID)
			router.Post("/todos/create", todoHandler.createTodo)
			router.Put("/todos/update/{id}", todoHandler.updateTodo)
			router.Patch("/todos/{id}/complete", todoHandler.toggleComplete)
			router.Delete("/todos/delete/{id}", todoHandler.deleteTodo)

			// Todo Details Routes
			router.Get("/todos/tododetails", todoTodoDetailsHandler.getTodoDetails)
			router.Get("/todos/tododetails/{id}", todoTodoDetailsHandler.getTodoDetailsByID)
			router.Get("/todos/tododetails/todoid/{todo_id}", todoTodoDetailsHandler.getTodoDetailsByTodoId)
			router.Post("/todos/{todo_id}/details", todoTodoDetailsHandler.createTodoDetails)
			router.Delete("/todos/tododetails/delete/{id}", todoTodoDetailsHandler.deleteTodoDetails)
		})

		router.Route("/v2", func(router chi.Router) {
			router.Get("/healthcheck", HealthCheck)
		})
	})

	return router

}
