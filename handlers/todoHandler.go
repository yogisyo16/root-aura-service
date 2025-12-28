package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-mongo-todos/services"
)

type TodoHandler struct {
	Service services.Todo
	services.TodoDetails
}

// Create Todo request structure
type CreateTodoRequest struct {
	Task      string `json:"task"`
	DateDue   string `json:"date_due"`
	Completed bool   `json:"completed"`
}

// Update Todo request structure
type UpdateTodoRequest struct {
	Task      string `json:"task"`
	DateDue   string `json:"date_due"`
	Completed bool   `json:"completed"`
}

// Generic response structure
func NewTodoHandler(service services.Todo) *TodoHandler {
	return &TodoHandler{
		Service: service,
	}
}

// Health Check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Msg:  "Health Check",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

// Logic to get all todos
func (h *TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.Service.GetAllTodos()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if todos == nil {
		todos = []services.Todo{}
	}

	// Response structure wrapper
	response := struct {
		Code int `json:"code"`
		Data struct {
			Items []services.Todo `json:"items"`
		} `json:"data"`
	}{
		Code: 200,
		Data: struct {
			Items []services.Todo `json:"items"`
		}{
			Items: todos,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

// Logic to get todo by id
func (h *TodoHandler) getTodoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	todo, err := h.Service.GetTodoById(id)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todo)
}

// Create todo
func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid request body",
			Code: 400,
		})
		return
	}

	// Parse the date_due string to time.Time
	dateDue, err := time.Parse(time.RFC3339, req.DateDue)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid date format. Use ISO 8601 format (e.g., 2025-11-20T14:30:00Z)",
			Code: 400,
		})
		return
	}

	// Create the Todo with parsed time.Time
	newTodo := services.Todo{
		Task:      req.Task,
		DateDue:   dateDue,
		Completed: req.Completed,
	}

	// Insert the new todo into the database
	err = h.Service.InsertTodo(newTodo)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Failed to create todo",
			Code: 500,
		})
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(Response{
		Msg:  "Successfully Created Todo",
		Code: 201,
	})
}

// Logic to update todo by id
func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid request body",
			Code: 400,
		})
		return
	}

	// Parse the date_due string to time.Time
	dateDue, err := time.Parse(time.RFC3339, req.DateDue)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid date format",
			Code: 400,
		})
		return
	}

	updateTodo := services.Todo{
		Task:      req.Task,
		DateDue:   dateDue, // Now it's time.Time
		Completed: req.Completed,
	}

	_, err = h.Service.UpdatedTodo(id, updateTodo)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			Msg:  err.Error(),
			Code: 500,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Response{
		Msg:  "Successfully Updated Todo",
		Code: 200,
	})
}

// Delete Todo it can be used only by id (On going to make delete all todos)
func (h *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Service.DeleteTodo(id)
	if err != nil {
		errorRes := Response{
			Msg:  "Error",
			Code: 304,
		}
		json.NewEncoder(w).Encode(errorRes)
		w.WriteHeader(errorRes.Code)
		return
	}

	res := Response{
		Msg:  "Succesfully Deleted Todo",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}
