package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Response struct {
	Msg  string
	Code int
}

func CreateRouter(todoHandler *TodoHandler, userHandler *UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CRSF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(router chi.Router) {
		router.Route("/v1", func(router chi.Router) {
			// User Routes
			router.Post("/users/create", userHandler.insertUser)
			router.Get("/users", userHandler.getAllUsers)
			router.Get("/users/{id}", userHandler.getUserByID)

			// Todo Routes
			router.Get("/healthcheck", HealthCheck)
			router.Get("/todos", todoHandler.getTodos)
			router.Get("/todos/{id}", todoHandler.getTodoByID)
			router.Post("/todos/create", todoHandler.createTodo)
			router.Put("/todos/update/{id}", todoHandler.updateTodo)
			router.Delete("/todos/delete/{id}", todoHandler.deleteTodo)
		})

		router.Route("/v2", func(router chi.Router) {
			router.Get("/healthcheck", HealthCheck)
		})
	})

	return router

}
