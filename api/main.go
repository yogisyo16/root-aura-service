package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yogisyo16/root-aura-service/db"
	"github.com/yogisyo16/root-aura-service/handlers"
	"github.com/yogisyo16/root-aura-service/services"
)

// Application groups the service layer for anything that wants a single
// handle on all of them. Not currently constructed/used anywhere below —
// main wires the individual services (todoService, userService,
// detailsService) straight into the handlers instead.
type Application struct {
	Models services.Models
}

// main wires the whole app together: connect to Mongo, construct the
// services, hand them to the handlers, build the router, start listening.
// See docs/HANDLERS.md and docs/SERVICES.md for what each layer does.
func main() {
	// 1. Connect to the database
	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	// 2. Initialize the services with the database client
	todoService := services.New(mongoClient)
	userService := services.User{} // Uses the same client set in services.New()
	detailsService := services.NewTodoDetailsService(mongoClient)

	// 3. Initialize the handlers with their respective services
	todoHandler := handlers.NewTodoHandler(todoService, detailsService) // Pass both services
	userHandler := handlers.NewUserHandler(userService)
	detailsHandler := handlers.NewTodoDetailsHandler(detailsService)

	// 4. Create the router and pass all handlers to it
	router := handlers.CreateRouter(todoHandler, userHandler, detailsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server is running on port %s\n", serverAddr)
	http.ListenAndServe(serverAddr, router)
}
